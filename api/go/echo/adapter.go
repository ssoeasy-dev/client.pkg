package ssoeasy

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	core "github.com/ssoeasy-dev/client.pkg/api/go/core"
	"github.com/ssoeasy-dev/client.pkg/api/go/core/client"
)

func GetUser(c echo.Context) *core.User {
	return core.GetUser(c.Request().Context())
}

func ExtractUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := core.ExtractUserFromRequest(c.Request())
			if user != nil {
				ctx := core.SetUser(c.Request().Context(), user)
				c.SetRequest(c.Request().WithContext(ctx))
			}
			return next(c)
		}
	}
}

func CheckPermission(client *client.Client, permissionID uuid.UUID) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, allowed, err := core.CheckPermissionFromRequest(c.Request().Context(), c.Request(), client, permissionID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "permission check failed")
			}
			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "forbidden")
			}
			ctx := core.SetUser(c.Request().Context(), user)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
