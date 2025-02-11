package acl

import (
	"context"
	"crypto/subtle"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/dennis-dko/go-toolkit/errorhandler"

	"github.com/casbin/casbin/v2"
	casbinmw "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const wildcard = "*"

var config *Config

type Config struct {
	Enabled     bool   `env:"ACL_ENABLED"`
	Username    string `env:"ACL_AUTH_USERNAME,unset"`
	Password    string `env:"ACL_AUTH_PASSWORD,unset"`
	AuthModel   string `env:"ACL_AUTH_MODEL"`
	PolicyModel string `env:"ACL_POLICY_MODEL"`
	enforcer    *casbin.Enforcer
}

// Provide provides configuration for acl
func (cfg *Config) Provide() error {
	enf, err := casbin.NewEnforcer(cfg.AuthModel, cfg.PolicyModel)
	if err != nil {
		return err
	}
	config = cfg
	config.enforcer = enf
	return nil
}

// AddUser adds a role for a user in the acl policy
func AddUser(ctx context.Context, userID string, roles []string) error {
	_, err := config.enforcer.AddRolesForUser(userID, roles)
	if err != nil {
		slog.ErrorContext(ctx, "failed to add user in acl policy", slog.Any("roles", roles), slog.String("error", err.Error()))
		return err
	}
	return nil
}

// DeleteUser deletes a role for a user in the acl policy
func DeleteUser(ctx context.Context, userID string) error {
	_, err := config.enforcer.DeleteUser(userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete user in acl policy", slog.String("error", err.Error()))
		return err
	}
	return nil
}

// GetPermissionsForUser gets all the permissions for a user
func GetPermissionsForUser(ctx context.Context, userID string) ([][]string, error) {
	perms, err := config.enforcer.GetImplicitPermissionsForUser(userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get acl permissions", slog.String("error", err.Error()))
		return nil, err
	}
	return perms, nil
}

// GetAuthorizedRoutes gets all the authorized routes
func GetAuthorizedRoutes() ([]string, error) {
	authRoutes, err := config.enforcer.GetAllObjects()
	if err != nil {
		return nil, err
	}
	return authRoutes, nil
}

// UseAuthEnforcer forces the acl on the routes
func UseAuthEnforcer(ctx context.Context, instance *echo.Echo) {
	if config.Enabled {
		instance.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
			Validator: func(username string, password string, c echo.Context) (bool, error) {
				if subtle.ConstantTimeCompare([]byte(username), []byte(config.Username)) == 1 &&
					subtle.ConstantTimeCompare([]byte(password), []byte(config.Password)) == 1 {
					slog.DebugContext(ctx, "Authentication is successfully for route", slog.String("route", c.Path()))
					return true, nil
				}
				return false, errorhandler.ErrAuthFailed
			},
			Skipper: func(c echo.Context) bool {
				perms, _ := GetPermissionsForUser(ctx, wildcard)
				for _, data := range perms {
					routePattern := data[1]
					if data[1] == wildcard {
						routePattern = fmt.Sprintf(".%s", wildcard)
					}
					matchRoute, err := regexp.MatchString(routePattern, c.Path())
					if err != nil {
						slog.ErrorContext(ctx, "error while matching the route", slog.String("error", err.Error()))
						return false
					}
					methodPattern := data[2]
					if data[2] == wildcard {
						methodPattern = fmt.Sprintf(".%s", wildcard)
					}
					matchMethod, err := regexp.MatchString(methodPattern, c.Request().Method)
					if err != nil {
						slog.ErrorContext(ctx, "error while matching the method", slog.String("error", err.Error()))
						return false
					}
					slog.DebugContext(ctx, "Permission check for route", slog.String("route", c.Path()),
						slog.String("method", c.Request().Method), slog.Bool("matchRoute", matchRoute),
						slog.Bool("matchMethod", matchMethod), slog.Any("permission", data))
					if matchRoute && matchMethod {
						slog.DebugContext(ctx, "Authentication skipped for route", slog.String("route", c.Path()))
						return true
					}
				}
				slog.DebugContext(ctx, "Authentication enforced for route", slog.String("route", c.Path()))
				return false
			},
		}))
		instance.Use(casbinmw.MiddlewareWithConfig(casbinmw.Config{
			Enforcer: config.enforcer,
			ErrorHandler: func(c echo.Context, internal error, proposedStatus int) error {
				slog.ErrorContext(ctx, "error while using the acl enforcer",
					slog.Int("status", proposedStatus), slog.String("error", internal.Error()),
				)
				return errorhandler.ErrPermFailed
			},
		}))
	} else {
		slog.InfoContext(ctx, "Authentication is disabled")
	}
}
