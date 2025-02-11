package logging

import (
	"context"
	"log/slog"
	"testing"

	"github.com/dennis-dko/go-toolkit/testhandler"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type LoggingTestSuite struct {
	suite.Suite
	ctx       context.Context
	instance  *echo.Echo
	logAttr   slog.Attr
	LogConfig *Config
}

func (l *LoggingTestSuite) SetupTest() {
	// Setup
	l.instance = echo.New()
	l.logAttr = slog.String("testKey", "testValue")
}

func (l *LoggingTestSuite) SetupSubTest() {
	// Sub setup
	l.ctx = testhandler.Ctx(false, false)
	l.LogConfig = &Config{}
}

func TestLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingTestSuite))
}

func (l *LoggingTestSuite) TestLogging() {

	l.Run("happy path - logging debug", func() {
		// Init
		l.LogConfig.LogLevelStr = "DEBUG"

		// Run
		err := l.LogConfig.Provide()

		// Assert
		l.NoError(err)
	})
	l.Run("happy path - logging info", func() {
		// Init
		l.LogConfig.LogLevelStr = "INFO"

		// Run
		err := l.LogConfig.Provide()

		// Assert
		l.NoError(err)
	})
	l.Run("happy path - logging warn", func() {
		// Init
		l.LogConfig.LogLevelStr = "WARN"

		// Run
		err := l.LogConfig.Provide()

		// Assert
		l.NoError(err)
	})
	l.Run("happy path - logging error", func() {
		// Init
		l.LogConfig.LogLevelStr = "ERROR"

		// Run
		err := l.LogConfig.Provide()

		// Assert
		l.NoError(err)
	})
	l.Run("should return an error by incorrect log level", func() {
		// Init
		l.LogConfig.LogLevelStr = "UNKNOWN"

		// Run
		err := l.LogConfig.Provide()

		// Assert
		l.Error(err)
		l.ErrorContains(err, "cannot provide UNKNOWN")
	})
}

func (l *LoggingTestSuite) TestAppendCtx() {

	l.Run("happy path - append log to context", func() {
		// Run
		logCtx := AppendCtx(l.ctx, l.logAttr)

		// Assert
		l.NotEmpty(logCtx.Value(slogFields))
	})
}

func (l *LoggingTestSuite) TestUseRequestLog() {

	l.Run("happy path - logging request", func() {
		// Init
		l.LogConfig.LogLevelStr = "INFO"

		// Run
		err := l.LogConfig.Provide()
		UseRequestLog(l.ctx, l.instance)

		// Assert
		l.NoError(err)
	})
}
