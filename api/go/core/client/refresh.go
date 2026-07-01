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

// Refresh выполняет HTTP-запрос /refresh и возвращает новые токены из заголовков ответа.
func (c *Client) Refresh(ctx context.Context, meta dto.Meta, tokens dto.Tokens) (dto.Tokens, error) {
	url := c.baseURL.JoinPath("/auth/refresh")

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url.String(), nil)
	if err != nil {
		return tokens, errors.NewWrapf(errors.ErrInternal, err, "refresh: build request")
	}

	if err := meta.ToHttpHeaders(req.Header); err != nil {
		return tokens, errors.NewWrapf(errors.ErrInvalidArgument, err, "refresh: set meta headers")
	}
	if err := tokens.ToHttpHeaders(req.Header); err != nil {
		return tokens, errors.NewWrapf(errors.ErrInvalidArgument, err, "refresh: set token headers")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return tokens, errors.NewWrapf(errors.ErrCanceled, err, "refresh: request canceled")
		}
		return tokens, errors.NewWrapf(errors.ErrBadGateway, err, "refresh: do request")
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return tokens, errors.NewWrapf(errors.ErrInternal, err, "refresh: read response body")
	}

	if res.StatusCode != http.StatusOK {
		kind := mapStatusCodeToKind(res.StatusCode)
		var errResp struct {
			Errors []string `json:"errors"`
		}
		_ = json.Unmarshal(resBody, &errResp)
		msg := "refresh failed"
		if len(errResp.Errors) > 0 {
			msg = errResp.Errors[0]
		}
		return tokens, errors.Newf(kind, "%s (status %d)", msg, res.StatusCode)
	}

	if tokens, err = dto.TokensFromHttpHeaders(res.Header); err != nil {
		return tokens, errors.NewWrapf(errors.ErrUnauthenticated, err, "refresh: invalid tokens in headers")
	}

	return tokens, nil
}
