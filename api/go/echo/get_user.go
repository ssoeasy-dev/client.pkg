package goecho

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// User содержит данные пользователя, извлечённые из JWT access-токена.
// Поля соответствуют структуре TokenClaims в auth.svc.
type User struct {
	UserID    uuid.UUID `json:"user_id"`
	ServiceID uuid.UUID `json:"service_id"`
	CompanyID uuid.UUID `json:"company_id"`
}

// rawClaims — внутренняя структура для json.Unmarshal из JWT payload.
type rawClaims struct {
	UserID    string  `json:"user_id"`
	ServiceID string  `json:"service_id"`
	CompanyID string  `json:"company_id"`
}

// ParseUser декодирует payload JWT access-токена и возвращает User.
//
// Подпись НЕ проверяется — это задача auth.svc. Функция только парсит
// уже доверенный токен, который прошёл через RequirePermission.
func ParseUser(token string) (*User, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("goecho: invalid JWT format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("goecho: base64 decode JWT payload: %w", err)
	}

	var raw rawClaims
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, fmt.Errorf("goecho: unmarshal JWT payload: %w", err)
	}

	userID, err := uuid.Parse(raw.UserID)
	if err != nil {
		return nil, fmt.Errorf("goecho: parse user_id: %w", err)
	}

	serviceID, err := uuid.Parse(raw.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("goecho: parse service_id: %w", err)
	}

	companyID, err := uuid.Parse(raw.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("goecho: parse company_id: %w", err)
	}

	return &User{
		UserID:    userID,
		ServiceID: serviceID,
		CompanyID: companyID,
	}, nil
}

// ─── Echo helpers ─────────────────────────────────────────────────────────────

const contextKeyUser = "goecho_user"

// GetUser достаёт *User из Echo-контекста.
// Работает только если перед этим был вызван миддлвар InjectUser (или RequirePermission с опцией WithInject).
// Возвращает nil если пользователь не был установлен.
func GetUser(c echo.Context) *User {
	u, _ := c.Get(contextKeyUser).(*User)
	return u
}

// InjectUser — миддлвар, который парсит токен и кладёт *User в Echo-контекст.
// Удобно ставить после RequirePermission, чтобы не парсить токен дважды.
//
//	e.Use(goecho.RequirePermission(client, permID), goecho.InjectUser())
//
// Или отдельно на роуты, где не нужна проверка разрешений, но нужен пользователь:
//
//	e.GET("/profile", profileHandler, goecho.InjectUser())
func InjectUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := extractBearer(c)
			if err != nil {
				// Токена нет — не падаем, просто не устанавливаем пользователя.
				return next(c)
			}

			user, err := ParseUser(token)
			if err != nil {
				// Токен есть, но невалидный — тоже не падаем здесь,
				// это забота RequirePermission.
				return next(c)
			}

			c.Set(contextKeyUser, user)
			return next(c)
		}
	}
}
