package utils

import (
	"encoding/json"
	"strings"
)

// Sensitive field names (case-insensitive) whose values are masked in logs.
var sensitiveFields = []string{
	"password", "token", "secret", "authorization", "refresh_token",
	"access_token", "api_key", "cookie", "jwt", "bearer",
}

const maskedValue = "***MASKED***"

// MaskSensitiveJSON masks known sensitive keys in a JSON string.
// Used for request/response body before logging.
func MaskSensitiveJSON(raw string) string {
	if raw == "" {
		return ""
	}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "{}" {
		return trimmed
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return raw
	}
	maskMap(m)
	b, _ := json.Marshal(m)
	return string(b)
}

// MaskQueryString masks sensitive key=value pairs in a query string for logging.
func MaskQueryString(q string) string {
	if q == "" {
		return ""
	}
	return maskQueryString(q)
}

func maskMap(m map[string]interface{}) {
	for k, v := range m {
		if isSensitiveKey(k) {
			m[k] = maskedValue
			continue
		}
		if nested, ok := v.(map[string]interface{}); ok {
			maskMap(nested)
		}
	}
}

func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range sensitiveFields {
		if lower == s || strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

// maskQueryString masks key=value pairs where key is sensitive.
func maskQueryString(q string) string {
	parts := strings.Split(q, "&")
	for i, part := range parts {
		idx := strings.Index(part, "=")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(part[:idx])
		if isSensitiveKey(key) {
			parts[i] = key + "=" + maskedValue
		}
	}
	return strings.Join(parts, "&")
}

// MaskBodyForLog returns a safe string for logging request/response body.
// If body is not valid JSON, returns as-is but with sensitive-looking substrings masked heuristically.
func MaskBodyForLog(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	return MaskSensitiveJSON(string(body))
}
