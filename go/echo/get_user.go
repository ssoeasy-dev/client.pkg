package goecho

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// User содержит данные пользователя, извлечённые из JWT access-токена.
// Поля соответствуют структуре TokenClaims в auth.svc.
type User struct {
	UserID    uuid.UUID `json:"user_id"`
	ServiceID uuid.UUID `json:"service_id"`
	CompanyID uuid.UUID `json:"company_id"`
	// Стандартные JWT claims.
	Subject   string    `json:"sub"`
	TokenID   string    `json:"jti"`
	ExpiresAt time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
}

// rawClaims — внутренняя структура для json.Unmarshal из JWT payload.
type rawClaims struct {
	UserID    string  `json:"user_id"`
	ServiceID string  `json:"service_id"`
	CompanyID string  `json:"company_id"`
	Type      string  `json:"type"`
	Sub       string  `json:"sub"`
	Jti       string  `json:"jti"`
	Exp       float64 `json:"exp"`
	Iat       float64 `json:"iat"`
}

// ParseUser декодирует payload JWT access-токена и возвращает User.
//
// Подпись НЕ проверяется — это задача auth.svc. Функция только парсит
// уже доверенный токен, который прошёл через RequirePermission.
func ParseUser(token string) (*User, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("authmw: invalid JWT format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("authmw: base64 decode JWT payload: %w", err)
	}

	var raw rawClaims
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, fmt.Errorf("authmw: unmarshal JWT payload: %w", err)
	}

	userID, err := uuid.Parse(raw.UserID)
	if err != nil {
		return nil, fmt.Errorf("authmw: parse user_id: %w", err)
	}

	serviceID, err := uuid.Parse(raw.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("authmw: parse service_id: %w", err)
	}

	companyID, err := uuid.Parse(raw.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("authmw: parse company_id: %w", err)
	}

	return &User{
		UserID:    userID,
		ServiceID: serviceID,
		CompanyID: companyID,
		Subject:   raw.Sub,
		TokenID:   raw.Jti,
		ExpiresAt: time.Unix(int64(raw.Exp), 0),
		IssuedAt:  time.Unix(int64(raw.Iat), 0),
	}, nil
}

// ─── Echo helpers ─────────────────────────────────────────────────────────────

const contextKeyUser = "authmw_user"

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
//	e.Use(authmw.RequirePermission(client, permID), authmw.InjectUser())
//
// Или отдельно на роуты, где не нужна проверка разрешений, но нужен пользователь:
//
//	e.GET("/profile", profileHandler, authmw.InjectUser())
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
