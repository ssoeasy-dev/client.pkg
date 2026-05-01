package ssoeasy

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// ExtractBearerTokenFromHeader извлекает токен из строки заголовка Authorization.
func ExtractBearerTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", NewError("missing Authorization header")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", NewError("invalid Authorization header format")
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", NewError("empty bearer token")
	}
	return token, nil
}

// ExtractBearerToken извлекает raw token из заголовка Authorization.
func ExtractBearerToken(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", NewError("missing Authorization header")
	}

	return ExtractBearerTokenFromHeader(h)
}

// DecodeJWTPayload декодирует payload из raw JWT токена.
func DecodeJWTPayload(token string) ([]byte, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, NewError("invalid JWT")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, NewError("invalid JWT: %w", err)
	}
	return payload, nil
}

// ExtractBearerTokenPayload извлекает и декодирует payload JWT.
func ExtractBearerTokenPayload(r *http.Request) ([]byte, error) {
	token, err := ExtractBearerToken(r)
	if err != nil {
		return nil, err
	}

	return DecodeJWTPayload(token)
}
