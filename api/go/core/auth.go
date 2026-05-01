package ssoeasy

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// ExtractUserFromRequest извлекает JWT из запроса, парсит пользователя
// и возвращает его. Если токен отсутствует или невалиден, возвращает nil.
// Ошибки парсинга игнорируются (возвращается nil).
func ExtractUserFromRequest(r *http.Request) *User {
	payload, err := ExtractBearerTokenPayload(r)
	if err != nil {
		return nil
	}
	user, err := ParseUser(payload)
	if err != nil {
		return nil
	}
	return user
}

// CheckPermissionFromRequest выполняет полную проверку:
// - извлекает токен из запроса
// - парсит пользователя
// - проверяет право через клиент
// Возвращает:
//   - user: распарсенный пользователь (nil при ошибке)
//   - allowed: true если доступ разрешен
//   - err: ошибка выполнения (не авторизационная, а сетевая/внутренняя)
func CheckPermissionFromRequest(
	ctx context.Context,
	r *http.Request,
	client Client,
	permissionID uuid.UUID,
) (user *User, allowed bool, err error) {
	token, err := ExtractBearerToken(r)
	if err != nil {
		return nil, false, nil // не авторизован, но это не ошибка выполнения
	}

	payload, err := ExtractBearerTokenPayload(r)
	if err != nil {
		return nil, false, nil
	}

	user, err = ParseUser(payload)
	if err != nil {
		return nil, false, nil
	}

	allowed, err = client.CheckPermission(ctx, token, permissionID)
	if err != nil {
		return user, false, err // ошибка выполнения (сеть, сервер)
	}
	return user, allowed, nil
}
