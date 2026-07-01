package echo

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ssoeasy-dev/client.pkg/api/go/core/v2"
	"github.com/ssoeasy-dev/client.pkg/api/go/core/v2/dto"
)

type Handler struct {
	service *core.Service
	cfg     CookieConfig
}

func NewHandler(cfg Config) (*Handler, error) {
	service, err := core.NewService(cfg.Client)
	if err != nil {
		return nil, err
	}

	return &Handler{
		service: service,
		cfg:     cfg.Cookie,
	}, nil
}

func (h *Handler) Authorize(cb core.Callback) echo.HandlerFunc {
	return func(c echo.Context) error {
		code, err := dto.ParseCode(c.Request())
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		meta, err := dto.MetaFromHttpHeaders(c.Request().Header)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		tokens, err := h.service.Authorize(c.Request().Context(), meta, code, cb)
		if err != nil {
			status := httpStatusFromError(err)
			return echo.NewHTTPError(status, err.Error())
		}

		if err := TokensFromCore(tokens).ToHttpCookie(c.Response().Writer, h.cfg); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "authorized"})
	}
}

func (h *Handler) Check(permissionID uuid.UUID, cb core.Callback) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			meta, err := dto.MetaFromHttpHeaders(c.Request().Header)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			tokens, err := TokensFromHttpCookie(c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			innerCb := func(ctx context.Context, payload dto.Payload) error {
                ctx = context.WithValue(ctx, dto.TokenPayloadContextKey, payload)
                c.SetRequest(c.Request().WithContext(ctx))
                c.Set(dto.TokenPayloadContextKey, payload)

                if cb != nil {
                    return cb(ctx, payload)
                }
                return nil
            }

			if err := h.service.Check(c.Request().Context(), meta, tokens, permissionID, innerCb); err != nil {
				status := httpStatusFromError(err)
				return echo.NewHTTPError(status, err.Error())
			}

			return next(c)
		}
	}
}

// Logout завершает сессию.
func (h *Handler) Logout(cb core.Callback) echo.HandlerFunc {
	return func(c echo.Context) error {
		meta, err := dto.MetaFromHttpHeaders(c.Request().Header)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		tokens, err := TokensFromHttpCookie(c.Request())
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		if err := h.service.Logout(c.Request().Context(), meta, tokens, cb); err != nil {
			status := httpStatusFromError(err)
			return echo.NewHTTPError(status, err.Error())
		}

		http.SetCookie(c.Response().Writer, &http.Cookie{
			Name:     dto.RefreshHeader,
			Value:    "",
			Path:     h.cfg.Path,
			Domain:   h.cfg.Domain,
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})

		return c.JSON(http.StatusOK, map[string]string{"status": "logged out"})
	}
}

// Me возвращает информацию о пользователе.
func (h *Handler) Me(cb core.Callback) echo.HandlerFunc {
	return func(c echo.Context) error {
		meta, err := dto.MetaFromHttpHeaders(c.Request().Header)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		tokens, err := TokensFromHttpCookie(c.Request())
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		user, err := h.service.Me(c.Request().Context(), meta, tokens, cb)
		if err != nil {
			status := httpStatusFromError(err)
			return echo.NewHTTPError(status, err.Error())
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Refresh обновляет токены.
func (h *Handler) Refresh(cb core.Callback) echo.HandlerFunc {
	return func(c echo.Context) error {
		meta, err := dto.MetaFromHttpHeaders(c.Request().Header)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		tokens, err := TokensFromHttpCookie(c.Request())
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		tokens, err = h.service.Refresh(c.Request().Context(), meta, tokens, cb)
		if err != nil {
			status := httpStatusFromError(err)
			return echo.NewHTTPError(status, err.Error())
		}

		if err := TokensFromCore(tokens).ToHttpCookie(c.Response().Writer, h.cfg); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "refreshed"})
	}
}
