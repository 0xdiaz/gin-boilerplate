package config

import (
	"fmt"
	"strings"
)

// defaultDevCORSOrigins are used for CORS_ALLOWED_ORIGINS when unset in development.
var defaultDevCORSOrigins = []string{
	"http://localhost:3000",
	"http://127.0.0.1:3000",
	"http://localhost:5173",
}

// defaultDevTrustedProxies are used for TRUSTED_PROXIES when unset in development.
var defaultDevTrustedProxies = []string{"127.0.0.1", "::1"}

// parseCommaSeparated splits a comma-separated env value into a trimmed, non-empty slice.
func parseCommaSeparated(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// normalizeCORSOrigins trims trailing slashes so "http://a.com/" and "http://a.com"
// are treated as the same origin.
func normalizeCORSOrigins(origins []string) []string {
	normalized := make([]string, 0, len(origins))
	for _, o := range origins {
		normalized = append(normalized, strings.TrimSuffix(o, "/"))
	}
	return normalized
}

// validateCORSOrigins rejects wildcard or malformed origins, and requires at least
// one explicit origin in production.
func validateCORSOrigins(rawOrigins []string, isProd bool) error {
	if isProd && len(rawOrigins) == 0 {
		return fmt.Errorf("CORS_ALLOWED_ORIGINS must be set in production (comma-separated frontend origins, e.g. https://app.example.com)")
	}
	for _, o := range rawOrigins {
		if o == "*" {
			return fmt.Errorf("CORS_ALLOWED_ORIGINS must not contain wildcard '*'; list explicit origins instead")
		}
		if !strings.HasPrefix(o, "http://") && !strings.HasPrefix(o, "https://") {
			return fmt.Errorf("CORS_ALLOWED_ORIGINS entry %q must start with http:// or https://", o)
		}
	}
	return nil
}

// resolveCORSOrigins normalizes explicitly configured origins, falling back to
// safe localhost defaults in development when none are set. In non-development
// environments an empty list means no cross-origin requests are allowed.
func resolveCORSOrigins(rawOrigins []string, isDev bool) []string {
	if len(rawOrigins) == 0 {
		if isDev {
			return append([]string(nil), defaultDevCORSOrigins...)
		}
		return nil
	}
	return normalizeCORSOrigins(rawOrigins)
}

// resolveTrustedProxies falls back to loopback-only defaults in development when
// TRUSTED_PROXIES is unset. In non-development environments an empty list means
// no proxy is trusted (Gin will use the raw remote address).
func resolveTrustedProxies(raw string, isDev bool) []string {
	proxies := parseCommaSeparated(raw)
	if len(proxies) == 0 && isDev {
		return append([]string(nil), defaultDevTrustedProxies...)
	}
	return proxies
}
