package config

import (
	"fmt"

	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
)

type ServerConfiguration struct {
	Host      string
	Port      string
	JWTSecret string
	Debug     bool
	// AllowedHosts is reserved for future Host-header validation; it is NOT
	// used for CORS (see CORSAllowedOrigins) or reverse proxy trust (see
	// TrustedProxies).
	AllowedHosts           string
	LimitCountPerRequest   int64
	Timezone               string // IANA name, e.g. UTC, Asia/Jakarta
	RequestTimeoutSeconds  int    // http.Server read/write timeout
	ShutdownTimeoutSeconds int    // graceful shutdown max wait
	RefreshTokenTTLDays    int    // refresh token lifetime before it must be re-issued via login
	AccessTokenTTLMinutes  int    // JWT access token lifetime; kept short since it cannot be revoked
	// Rate limit (used by middlewares)
	RateLimitRPS     int
	RateLimitBurst   int
	RateLimitUseUser bool
	// CORSAllowedOrigins is the allowlist of exact frontend origins permitted
	// to make cross-origin requests. Never contains "*".
	CORSAllowedOrigins []string
	// TrustedProxies is the list of IPs/CIDRs trusted to set X-Forwarded-For,
	// passed to gin.Engine.SetTrustedProxies(). Unrelated to CORS.
	TrustedProxies []string
}

// ServerConfig returns the server address (host:port) from loaded config.
// Call after SetupConfig(). Returns default "0.0.0.0:8000" if config not loaded.
func ServerConfig() string {
	if c := Get(); c != nil {
		addr := fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
		logger.Infof("server listening at %s", addr)
		return addr
	}
	logger.Infof("server listening at 0.0.0.0:8000 (config not loaded)")
	return "0.0.0.0:8000"
}
