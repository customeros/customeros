package utils

import (
	"golang.org/x/net/publicsuffix"
	"net/url"
	"strings"
)

// GetDomainPrefix returns everything before the first dot in the domain
func GetDomainPrefix(domain string) string {
	// Split the domain by dots
	parts := strings.Split(domain, ".")
	// Return the first part
	if len(parts) > 0 {
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
	etld, im := publicsuffix.PublicSuffix(input)
	var validtld = false
	if im { // ICANN managed
		validtld = true
	} else if strings.IndexByte(etld, '.') >= 0 { // privately managed
		validtld = true
	}
	return validtld
}
