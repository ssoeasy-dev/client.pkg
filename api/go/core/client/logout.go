package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ssoeasy-dev/client.pkg/api/go/core/v2/dto"
	"github.com/ssoeasy-dev/pkg/errors"
)

func (c *Client) Logout(ctx context.Context, meta dto.Meta, tokens dto.Tokens) error {
	url := c.baseURL.JoinPath("/auth/logout")

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), nil)
	if err != nil {
		return errors.NewWrapf(errors.ErrInternal, err, "logout: build request")
	}

	if err := meta.ToHttpHeaders(req.Header); err != nil {
		return errors.NewWrapf(errors.ErrInvalidArgument, err, "logout: set meta headers")
	}
	if err := tokens.ToHttpHeaders(req.Header); err != nil {
		return errors.NewWrapf(errors.ErrInvalidArgument, err, "logout: set token headers")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.NewWrapf(errors.ErrCanceled, err, "logout: request canceled")
		}
		return errors.NewWrapf(errors.ErrBadGateway, err, "logout: do request")
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.NewWrapf(errors.ErrInternal, err, "logout: read response body")
	}

	if res.StatusCode != http.StatusOK {
		kind := mapStatusCodeToKind(res.StatusCode)
		var errResp struct {
			Errors []string `json:"errors"`
		}
		_ = json.Unmarshal(resBody, &errResp)
		msg := "logout failed"
		if len(errResp.Errors) > 0 {
			msg = errResp.Errors[0]
		}
		return errors.Newf(kind, "%s (status %d)", msg, res.StatusCode)
	}

	return nil
}
