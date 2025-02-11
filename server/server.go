package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dennis-dko/go-toolkit/tracing"

	"github.com/dennis-dko/go-toolkit/acl"
	"github.com/dennis-dko/go-toolkit/logging"
	"github.com/dennis-dko/go-toolkit/recoverhandler"
	"github.com/dennis-dko/go-toolkit/secure"

	"github.com/labstack/echo/v4"
)

// signalNotify wraps signal.Notify for test purposes
var signalNotify = func(c chan os.Signal, sig ...os.Signal) {
	signal.Notify(c, sig...)
}

type Config struct {
	Name                    string        `env:"NAME"`
	Host                    string        `env:"HOST,notEmpty"`
	Port                    int           `env:"PORT,notEmpty"`
	IsProduction            bool          `env:"PRODUCTION"`
	GracefulShutDownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"60s"`
	Logging                 logging.Config
	Tracing                 tracing.Config
	Recover                 recoverhandler.Config
	Secure                  secure.Config
	Acl                     acl.Config
}

type Server struct {
	Name     string
	Echo     *echo.Echo
	Context  context.Context
	Config   *Config
	EnvFiles []string
}

// New creates a new server instance
func New(ctx context.Context, cfg *Config, envFiles []string) *Server {
	echoInit := echo.New()
	if cfg.IsProduction {
		echoInit.HideBanner = true
		echoInit.HidePort = true
	}
	return &Server{
		Name:     cfg.Name,
		Echo:     echoInit,
		Context:  ctx,
		Config:   cfg,
		EnvFiles: envFiles,
	}
}

// Start starting the server
func (server *Server) Start() {
	go func() {
		if len(server.EnvFiles) > 0 {
			slog.InfoContext(server.Context, "Env files will be used for this server", slog.String("serverName", server.Name), slog.String("envFiles", strings.Join(server.EnvFiles, ",")))
		}
		err := server.Echo.Start(
			fmt.Sprintf("%s:%d", server.Config.Host, server.Config.Port),
		)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(server.Context, "error while starting the server, terminating", slog.String("serverName", server.Name), slog.String("error", err.Error()))
		}
	}()
	// Wait for interrupt signal to gracefully shut down the server with a specified timeout.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := CloseProcess()
	<-quit
	cancelCtx, cancel := context.WithTimeout(server.Context, server.Config.GracefulShutDownTimeout)
	defer cancel()
	if err := server.Echo.Shutdown(cancelCtx); err != nil {
		slog.ErrorContext(server.Context, "failed to gracefully shutdown the server", slog.String("serverName", server.Name), slog.String("error", err.Error()))
		os.Exit(1)
	}
}

// CloseProcess closes the os signal process
func CloseProcess() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signalNotify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	return quit
}

// IsProcessClosed checks if the given os signal process is closed
func IsProcessClosed(ch <-chan os.Signal) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
