package client

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/ssoeasy-dev/client.pkg/api/go/core"
)

// CheckPermission выполняет GET /api/v1/permission/check и возвращает:
//   - true  — доступ разрешён (HTTP 200)
//   - false — доступ запрещён (HTTP 403)
//   - error — сетевая / протокольная ошибка
func (c *Client) CheckPermission(
	ctx context.Context,
	token string,
	permissionID uuid.UUID,
) (bool, error) {
	endpoint := c.baseURL + endpointCheckPermission

	q := url.Values{}
	q.Set("permissionId", permissionID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return false, ssoeasy.NewError("check permission build request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, ssoeasy.NewError("check permission request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusForbidden:
		return false, nil
	case http.StatusUnauthorized:
		return false, ssoeasy.NewError("unauthorized")
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return false, ssoeasy.NewError("check permission unexpected status %d : %s", resp.StatusCode, body)
	}
}
