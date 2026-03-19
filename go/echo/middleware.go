// Package authmw предоставляет Echo-миддлвар для проверки ABAC-разрешений
// через HTTP-запрос к auth.api.
//
// Использование:
//
//	client := authmw.NewClient("https://auth.api.example.com")
//
//	e := echo.New()
//	e.GET("/orders", authmw.RequirePermission(client, "f47ac10b-..."), handler)
package goecho

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)



// ─── Echo middleware ──────────────────────────────────────────────────────────

// EchoMiddleware возвращает Echo-миддлвар, который:
//  1. Извлекает Bearer-токен из заголовка Authorization.
//  2. Опционально читает company_id из заголовка X-Company-Id.
//  3. Вызывает GET /api/v1/permission/check на auth.api.
//  4. Возвращает 401 / 403 при отказе, или пропускает запрос дальше.
//
// permissionID — UUID разрешения из таблицы permissions в auth.svc.
func EchoMiddleware(client *Client, permissionID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := extractBearer(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid Authorization header")
			}

			companyID := c.Request().Header.Get("X-Company-Id")

			allowed, err := client.CheckPermission(c.Request().Context(), token, permissionID, companyID)
			if err != nil {
				if err == ErrUnauthorized {
					return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
				}
				// Сетевая ошибка — не пропускаем, возвращаем 503.
				return echo.NewHTTPError(http.StatusServiceUnavailable, "permission check unavailable")
			}

			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}

			return next(c)
		}
	}
}

// extractBearer достаёт токен из "Authorization: Bearer <token>".
func extractBearer(c echo.Context) (string, error) {
	h := c.Request().Header.Get("Authorization")
	if h == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", fmt.Errorf("invalid Authorization header format")
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", fmt.Errorf("empty bearer token")
	}
	return token, nil
}
