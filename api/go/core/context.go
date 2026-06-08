package ssoeasy

import "context"

type contextKey string

const userContextKey contextKey = "ssoeasy_user"

// SetUser добавляет пользователя в контекст.
func SetUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUser извлекает пользователя из контекста.
func GetUser(ctx context.Context) *User {
	u, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil
	}
	return u
}
