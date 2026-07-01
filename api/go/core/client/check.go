package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ssoeasy-dev/client.pkg/api/go/core/v2/dto"
	"github.com/ssoeasy-dev/pkg/errors"
)

// Check выполняет GET /api/v1/permission/check и возвращает:
//   - error — если доступ запрещён или произошла ошибка.
func (c *Client) Check(ctx context.Context, meta dto.Meta, tokens dto.Tokens, permissionID uuid.UUID) error {
	path := "/permission/check/" + permissionID.String()
	url := c.baseURL.JoinPath(path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return errors.NewWrapf(errors.ErrInternal, err, "check: build request")
	}

	if err = meta.ToHttpHeaders(req.Header); err != nil {
		return errors.NewWrapf(errors.ErrInvalidArgument, err, "check: set meta headers")
	}
	if err = tokens.ToHttpHeaders(req.Header); err != nil {
		return errors.NewWrapf(errors.ErrInvalidArgument, err, "check: set token headers")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.NewWrapf(errors.ErrCanceled, err, "check: request canceled")
		}
		return errors.NewWrapf(errors.ErrBadGateway, err, "check: do request")
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.NewWrapf(errors.ErrInternal, err, "check: read response body")
	}

	if res.StatusCode != http.StatusOK {
		kind := mapStatusCodeToKind(res.StatusCode)
		var errResp struct {
			Errors []string `json:"errors"`
		}
		_ = json.Unmarshal(resBody, &errResp)
		msg := "permission check failed"
		if len(errResp.Errors) > 0 {
			msg = errResp.Errors[0]
		}
		return errors.Newf(kind, "%s (status %d)", msg, res.StatusCode)
	}

	return nil
}
