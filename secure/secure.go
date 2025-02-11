package secure

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/dennis-dko/go-toolkit/errorhandler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

var config *Config

type Config struct {
	Enabled               bool          `env:"SECURE_ENABLED"`
	XSSProtection         string        `env:"SECURE_HEADER_XSS" envDefault:"1; mode=block"`
	ContentTypeNosniff    string        `env:"SECURE_HEADER_NO_SNIFF" envDefault:"nosniff"`
	XFrameOptions         string        `env:"SECURE_HEADER_XFRAME" envDefault:"SAMEORIGIN"`
	HSTSMaxAge            int           `env:"SECURE_HEADER_MAX_AGE" envDefault:"3600"`
	ContentSecurityPolicy string        `env:"SECURE_HEADER_CSP" envDefault:"default-src 'self'"`
	AllowHeaders          []string      `env:"SECURE_CORS_ALLOW_HEADERS"`
	AllowMethods          []string      `env:"SECURE_CORS_ALLOW_METHODS"`
	AllowOrigins          []string      `env:"SECURE_CORS_ALLOW_ORIGINS" envDefault:"*"`
	AllowCredentials      bool          `env:"SECURE_CORS_ALLOW_CREDENTIALS"`
	RateLimit             float64       `env:"SECURE_RATE_LIMIT" envDefault:"10"`
	Burst                 int           `env:"SECURE_RATE_BURST" envDefault:"30"`
	ExpiresIn             time.Duration `env:"SECURE_RATE_EXPIRES_IN" envDefault:"3m"`
	TokenLength           uint8         `env:"SECURE_CSRF_TOKEN_LENGTH" envDefault:"32"`
	TokenLookup           string        `env:"SECURE_CSRF_TOKEN_HEADER" envDefault:"X-CSRF-Token"`
	CookieName            string        `env:"SECURE_CSRF_COOKIE_NAME" envDefault:"_csrf"`
	CookieMaxAge          int           `env:"SECURE_CSRF_COOKIE_MAX_AGE" envDefault:"86400"`
	CookieSecure          bool          `env:"SECURE_CSRF_COOKIE_SECURE"`
}

// Provide provides configuration for secure
func (cfg *Config) Provide() {
	config = cfg
}

// UseSecure enables security rules
func UseSecure(ctx context.Context, instance *echo.Echo) {
	if config.Enabled {
		instance.Use(middleware.SecureWithConfig(middleware.SecureConfig{
			XSSProtection:         config.XSSProtection,
			ContentTypeNosniff:    config.ContentTypeNosniff,
			XFrameOptions:         config.XFrameOptions,
			HSTSMaxAge:            config.HSTSMaxAge,
			ContentSecurityPolicy: config.ContentSecurityPolicy,
		}))
		instance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowHeaders:     config.AllowHeaders,
			AllowMethods:     config.AllowMethods,
			AllowOrigins:     config.AllowOrigins,
			AllowCredentials: config.AllowCredentials,
		}))
		instance.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{
					Rate:      rate.Limit(config.RateLimit),
					Burst:     config.Burst,
					ExpiresIn: config.ExpiresIn,
				},
			),
			ErrorHandler: func(c echo.Context, err error) error {
				slog.ErrorContext(ctx, "error while using the rate limit",
					slog.String("error", err.Error()),
				)
				return errorhandler.ErrAuthFailed
			},
			DenyHandler: func(c echo.Context, identifier string, err error) error {
				slog.InfoContext(ctx, "Access denied while sending too many requests",
					slog.String("identifier", identifier), slog.String("error", err.Error()),
				)
				return errorhandler.ErrRequestsLimitExceeded
			},
		}))
		instance.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLength:  config.TokenLength,
			TokenLookup:  fmt.Sprintf("header:%s", config.TokenLookup),
			ContextKey:   strings.Trim(config.CookieName, "_"),
			CookieName:   config.CookieName,
			CookieMaxAge: config.CookieMaxAge,
		}))
	} else {
		slog.InfoContext(ctx, "Web secure is disabled")
	}
}
