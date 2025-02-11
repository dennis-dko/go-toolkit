package recoverhandler

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var config *Config

type Config struct {
	StackSize         int  `env:"RECOVER_STACK_SIZE" envDefault:"4096"` // 4 KB
	DisableStackAll   bool `env:"RECOVER_DISABLE_STACK_ALL"`
	DisablePrintStack bool `env:"RECOVER_DISABLE_PRINT_STACK"`
}

// Provide provides configuration for recover
func (cfg *Config) Provide() {
	config = cfg
}

// UseRecover recovers by panic
func UseRecover(ctx context.Context, instance *echo.Echo) {
	instance.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         config.StackSize,
		DisableStackAll:   config.DisableStackAll,
		DisablePrintStack: config.DisablePrintStack,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			if slog.Default().Enabled(ctx, slog.LevelDebug) {
				slog.DebugContext(ctx, "PANIC RECOVER",
					slog.String("error", err.Error()),
					slog.String("stack", string(stack)),
				)
			} else {
				slog.ErrorContext(ctx, "PANIC RECOVER",
					slog.String("error", err.Error()),
					slog.String("stack", string(stack)),
				)
			}
			return err
		},
	}))
}
