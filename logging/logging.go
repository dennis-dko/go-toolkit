package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/dennis-dko/go-toolkit/constant"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/greyxor/slogor"
)

const slogFields = "slog_fields"

type Config struct {
	LogLevelStr string `env:"LOG_LEVEL"`
	LogAsJson   bool   `env:"LOG_AS_JSON"`
	LogLevel    slog.Level
}

type ContextHandler struct {
	slog.Handler
}

// Provide provides configuration for logging
func (cfg *Config) Provide() error {
	switch strings.ToLower(cfg.LogLevelStr) {
	case "debug":
		cfg.LogLevel = slog.LevelDebug
	case "info":
		cfg.LogLevel = slog.LevelInfo
	case "warn":
		cfg.LogLevel = slog.LevelWarn
	case "error":
		cfg.LogLevel = slog.LevelError
	default:
		return fmt.Errorf("cannot provide %s", cfg.LogLevelStr)
	}
	var logger *slog.Logger
	if cfg.LogAsJson {
		jsonHandler := &ContextHandler{
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					AddSource:   true,
					Level:       cfg.LogLevel,
					ReplaceAttr: replaceMsgKey(),
				},
			),
		}
		logger = slog.New(jsonHandler)
	} else {
		textHandler := &ContextHandler{
			slogor.NewHandler(
				os.Stdout,
				slogor.ShowSource(),
				slogor.SetTimeFormat(time.Stamp),
				slogor.SetLevel(cfg.LogLevel),
			),
		}
		logger = slog.New(textHandler)
	}
	slog.SetDefault(logger)
	return nil
}

// Handle logs slog attributes
func (ch ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}
	return ch.Handler.Handle(ctx, r)
}

// AppendCtx appends slog attributes to context
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}
	var v []slog.Attr
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}

// UseRequestLog logs request meta data
func UseRequestLog(ctx context.Context, instance *echo.Echo) {
	instance.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID: true,
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(ctx, slog.LevelInfo, "REQUEST",
					slog.String("id", v.RequestID),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(ctx, slog.LevelError, "REQUEST_ERROR",
					slog.String("id", v.RequestID),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("error", v.Error.Error()),
				)
			}
			return nil
		},
	}))
}

// UseBodyDump logs request and response bodies (only if debug level is enabled)
func UseBodyDump(ctx context.Context, instance *echo.Echo, skipUrls ...string) {
	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		instance.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
			Skipper: func(c echo.Context) bool {
				if slices.Contains(skipUrls, c.Request().URL.Path) {
					return true
				}
				return false
			},
			Handler: func(c echo.Context, reqBody, resBody []byte) {
				slog.LogAttrs(ctx, slog.LevelDebug, "BODY_DUMP",
					slog.String("request", string(reqBody)),
					slog.String("response", string(resBody)),
				)
			},
		}))
	} else {
		slog.InfoContext(ctx, "Body dump is disabled")
	}
}

func replaceMsgKey() func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.MessageKey {
			return slog.Attr{
				Key:   constant.MessageLogKey,
				Value: a.Value,
			}
		}
		return a
	}
}
