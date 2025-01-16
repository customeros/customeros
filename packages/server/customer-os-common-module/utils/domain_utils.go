package utils

import (
	"golang.org/x/net/publicsuffix"
	"net/url"
	"strings"
)

var domainExceptions = []string{
	"nhs.uk",
	"gov.uk",
	"ac.uk",
	"mod.uk",
	"parliament.uk",
	"police.uk",
	"ltd.uk",
	"plc.uk",
	"me.uk",
	"sch.uk",
}

func isDomainException(hostname string) bool {
	for _, ex := range domainExceptions {
		if hostname == ex {
			return true
		}
		if strings.HasSuffix(hostname, "."+ex) {
			return true
		}
	}
	return false
}

func tldPlusOneException(hostname string) string {
	for _, ex := range domainExceptions {
		// If exactly the exception
		if hostname == ex {
			return ex
		}
		// If it ends with our exception (e.g. "my.sub.nhs.uk" ends with ".nhs.uk")
		if strings.HasSuffix(hostname, "."+ex) {
			// Remove the ".nhs.uk" part
			label := strings.TrimSuffix(hostname, "."+ex)
			parts := strings.Split(label, ".")
			// If there's only one label before the exception, e.g., "sub.nhs.uk"
			if len(parts) == 1 {
				return parts[0] + "." + ex // → "sub.nhs.uk"
			}
			// If multiple labels, e.g. "my.sub.nhs.uk", keep the last label + ex
			lastPart := parts[len(parts)-1]
			return lastPart + "." + ex // → "sub.nhs.uk"
		}
	}
	return ""
}

// GetDomainWithoutTLD returns everything before the last dot in the domain
func GetDomainWithoutTLD(domain string) string {
	// Split the domain by dots
	parts := strings.Split(domain, ".")
	// Return all but last part
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ".")
	} else if len(parts) == 1 {
		return parts[0]
	}
	return domain
}

func ExtractDomain(input string) string {
	if !strings.Contains(input, ".") {
		return ""
	}

	hostname := extractHostname(strings.TrimSpace(strings.ToLower(input)))

	domain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		// PSL failed; check our exception fallback
		fallbackDomain := tldPlusOneException(hostname)
		if fallbackDomain != "" {
			return fallbackDomain
		}
		return ""
	}

	if IsValidTLD(domain) {
		return domain
	}
	return ""
}

func extractHostname(inputURL string) string {
	// Prepend "http://" if the URL doesn't start with a scheme
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "http://" + inputURL
	}

	// Parse the URL
	u, err := url.Parse(inputURL)
	if err != nil {
		return ""
	}

	// Extract and return the hostname (domain)
	hostname := u.Hostname()

	// Remove "www." if it exists
	if strings.HasPrefix(hostname, "www.") {
		hostname = hostname[4:] // Remove the first 4 characters ("www.")
	}

	return strings.ToLower(hostname)
}

func IsValidTLD(input string) bool {
	// If the entire domain is in or ends with an exception, consider it valid
	if isDomainException(input) {
		return true
	}

	etld, icannManaged := publicsuffix.PublicSuffix(input)

	// PSL recognized the domain
	if icannManaged {
		return true
	}
	// PSL says it's privately managed (i.e. something like 'co.uk' or 'appspot.com')
	if strings.Contains(etld, ".") {
		return true
	}
	// Otherwise, we consider it invalid
	return false
}

func IsValidDomain(input string) bool {
	// Quick reject if input looks like a URL with paths or queries
	if strings.ContainsAny(input, "/?&") {
		return false
	}
	// Reject if starting with "www."
	if strings.HasPrefix(input, "www.") {
		return false
	}

	_, err := publicsuffix.EffectiveTLDPlusOne(input)
	if err == nil {
		// PSL was happy
		return true
	}
	// PSL errored; check the exception
	return isDomainException(input)
}
