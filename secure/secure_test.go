package secure

import (
	"context"
	"testing"
	"time"

	"github.com/dennis-dko/go-toolkit/testhandler"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type SecureTestSuite struct {
	suite.Suite
	ctx      context.Context
	instance *echo.Echo
	config   Config
}

func (s *SecureTestSuite) SetupTest() {
	// Setup
	s.instance = echo.New()
	s.config = Config{
		Enabled:               true,
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
		AllowHeaders:          []string{},
		AllowMethods:          []string{},
		AllowOrigins:          []string{"*"},
		AllowCredentials:      false,
		RateLimit:             10,
		Burst:                 30,
		ExpiresIn:             3 * time.Minute,
		TokenLength:           32,
		TokenLookup:           echo.HeaderXCSRFToken,
		CookieName:            "_csrf",
		CookieMaxAge:          86400,
		CookieSecure:          false,
	}
}

func (s *SecureTestSuite) SetupSubTest() {
	// Sub setup
	s.ctx = testhandler.Ctx(false, false)
}

func TestSecureTestSuite(t *testing.T) {
	suite.Run(t, new(SecureTestSuite))
}

func (s *SecureTestSuite) TestUseEnforcer() {

	s.Run("happy path - use secure", func() {
		// Run
		s.config.Provide()
		UseSecure(s.ctx, s.instance)

		// Assert
		s.Equal(s.config, *config)
	})
}
