package goecho

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client — HTTP-клиент к auth.api.
// Создайте один раз и переиспользуйте во всём приложении.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// ClientOption позволяет настроить Client.
type clientOption func(*Client)

// NewClient создаёт Client.
//   - baseURL — корень auth.api, например "https://auth.ssoeasy.ru"
func NewClient(baseURL string, opts ...clientOption) *Client {
	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// CheckPermission выполняет GET /api/v1/permission/check и возвращает:
//   - true  — доступ разрешён (HTTP 200)
//   - false — доступ запрещён (HTTP 403)
//   - error — сетевая / протокольная ошибка
func (c *Client) CheckPermission(
	ctx context.Context,
	token string,
	permissionID string,
	companyID string, // пустая строка — не передавать
) (bool, error) {
	endpoint := c.baseURL + "/api/v1/permission/check"

	q := url.Values{}
	q.Set("permissionId", permissionID)
	if companyID != "" {
		q.Set("companyId", companyID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return false, fmt.Errorf("goecho: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("goecho: request to auth.api: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Читаем тело только при ошибках, чтобы включить в сообщение.
	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusForbidden:
		return false, nil
	case http.StatusUnauthorized:
		return false, ErrUnauthorized
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return false, fmt.Errorf("goecho: unexpected status %d: %s", resp.StatusCode, body)
	}
}

// ErrUnauthorized возвращается когда auth.api ответил 401.
var ErrUnauthorized = fmt.Errorf("goecho: unauthorized")
