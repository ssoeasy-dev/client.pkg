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

func (c *Client) Me(ctx context.Context, meta dto.Meta, tokens dto.Tokens) (dto.User, error) {
	var user dto.User
	url := c.baseURL.JoinPath("/auth/me")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return user, errors.NewWrapf(errors.ErrInternal, err, "me: build request")
	}

	if err := meta.ToHttpHeaders(req.Header); err != nil {
		return user, errors.NewWrapf(errors.ErrInvalidArgument, err, "me: set meta headers")
	}
	if err := tokens.ToHttpHeaders(req.Header); err != nil {
		return user, errors.NewWrapf(errors.ErrInvalidArgument, err, "me: set token headers")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return user, errors.NewWrapf(errors.ErrCanceled, err, "me: request canceled")
		}
		return user, errors.NewWrapf(errors.ErrBadGateway, err, "me: do request")
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return user, errors.NewWrapf(errors.ErrInternal, err, "me: read response body")
	}

	if res.StatusCode != http.StatusOK {
		kind := mapStatusCodeToKind(res.StatusCode)
		var errResp struct {
			Errors []string `json:"errors"`
		}
		_ = json.Unmarshal(resBody, &errResp)
		msg := "me request failed"
		if len(errResp.Errors) > 0 {
			msg = errResp.Errors[0]
		}
		return user, errors.Newf(kind, "%s (status %d)", msg, res.StatusCode)
	}

	if err = json.Unmarshal(resBody, &user); err != nil {
		return user, errors.NewWrapf(errors.ErrInternal, err, "me: unmarshal response")
	}

	return user, nil
}
