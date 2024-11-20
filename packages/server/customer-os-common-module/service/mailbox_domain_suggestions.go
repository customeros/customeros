package service

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

// checkDomain checks if a domain is likely available using multiple methods
func (s *mailboxService) IsDomainAvailable(ctx context.Context, domain string) (ok, available bool) {
	span, ctx := s.initializeTracing(ctx, "MailboxService.IsDomainAvailable")
	span.LogFields(log.String("domain", domain))
	defer span.Finish()

	// First try DNS lookup
	dnsOk, dnsCheck := s.dnsCheck(domain, span)
	if dnsCheck {
		return true, false
	}

	whoOk, whoCheck := s.checkWhois(domain, span)

	if !dnsOk && !whoOk {
		return false, false
	}

	if whoCheck {
		return true, false
	}

	return true, true
}

func (s *mailboxService) RecommendOutboundDomains(ctx context.Context, domainRoot string, count int) []string {
	span, ctx := s.initializeTracing(ctx, "MailboxService.RecommendOutboundDomains")
	span.LogFields(log.String("domainRoot", domainRoot))
	defer span.Finish()

	var results []string
	c := 0
	pre := getDomainPrefix()
	suf := getDomainSuffix()

	for _, prefix := range pre {
		if c == count {
			return results
		}
		newDomain := fmt.Sprintf("%s%s%s", prefix, domainRoot, ".com")
		_, avaliable := s.IsDomainAvailable(ctx, newDomain)
		if avaliable {
			c++
			results = append(results, newDomain)
		}
	}

	for _, suffix := range suf {
		if c == count {
			return results
		}
		newDomain := fmt.Sprintf("%s%s%s", domainRoot, suffix, ".com")
		_, available := s.IsDomainAvailable(ctx, newDomain)
		if available {
			c++
			results = append(results, newDomain)
		}
	}

	for _, prefix := range pre {
		for _, suffix := range suf {
			if c == count {
				return results
			}
			newDomain := fmt.Sprintf("%s%s%s%s", prefix, domainRoot, suffix, ".com")
			_, available := s.IsDomainAvailable(ctx, newDomain)
			if available {
				c++
				results = append(results, newDomain)
			}
		}
	}

	return results
}

func (s *mailboxService) dnsCheck(domain string, span opentracing.Span) (ok bool, exists bool) {
	ips, err := net.LookupIP(domain)
	if len(ips) > 0 {
		return true, true
	}

	if err != nil && strings.Contains(err.Error(), "no such host") {
		return true, false
	}

	if err != nil {
		tracing.TraceErr(span, err)
		return false, false
	}

	return true, false
}

// checkWhois runs a whois query and analyzes the output
func (s *mailboxService) checkWhois(domain string, span opentracing.Span) (ok, exists bool) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create the command with context
	cmd := exec.CommandContext(ctx, "whois", domain)

	// Create a pipe for stdout
	output, err := cmd.Output()

	// Check if the context deadline exceeded
	if ctx.Err() == context.DeadlineExceeded {
		return false, false
	}

	if err != nil {
		tracing.TraceErr(span, err)
		// Check if it's an exit error (whois sometimes exits with status 1)
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Still process the output if we got any
			if len(exitErr.Stderr) > 0 {
				output = exitErr.Stderr
			}
		} else {
			return false, false
		}
	}

	response := strings.ToLower(string(output))

	// Common phrases indicating domain is available
	availablePhrases := []string{
		"no match for",
		"not found",
		"no entries found",
		"no data found",
		"domain not found",
		"status: free",
		"status: available",
	}

	// Common phrases indicating domain is taken
	takenPhrases := []string{
		"domain name:",
		"registrar:",
		"creation date:",
		"registered on:",
		"status: active",
	}

	// Check for availability indicators
	for _, phrase := range availablePhrases {
		if strings.Contains(response, phrase) {
			return true, false
		}
	}

	// Check for registered indicators
	for _, phrase := range takenPhrases {
		if strings.Contains(response, phrase) {
			return true, true
		}
	}

	return true, false
}

func getDomainSuffix() []string {
	suffix := []string{
		"ai",
		"hq",
		"io",
		"ly",
		"app",
		"dev",
		"api",
		"hub",
		"now",
		"tech",
		"labs",
		"zone",
		"team",
		"tools",
		"cloud",
		"software",
		"platform",
	}
	return suffix
}

func getDomainPrefix() []string {
	prefix := []string{
		"go",
		"by",
		"at",
		"get",
		"try",
		"try",
		"use",
		"run",
		"meet",
		"join",
		"from",
		"with",
	}
	return prefix
}
