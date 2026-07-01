package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ssoeasy-dev/pkg/errors"

	"github.com/ssoeasy-dev/client.pkg/api/go/core/v2/dto"
)

// Authorize выполняет первоначальный обмен кода авторизации на пару токенов.
// POST /api/v1/auth/authorize/{serviceID}
func (c *Client) Authorize(ctx context.Context, meta dto.Meta, code dto.Code) (dto.Tokens, error) {
	var tokens dto.Tokens
	path := "/auth/authorize/" + c.serviceID.String()
	url := c.baseURL.JoinPath(path)

	body, err := json.Marshal(code)
	if err != nil {
		return tokens, errors.NewWrap(errors.ErrInvalidArgument, err, "authorize marshal request error")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(body))
	if err != nil {
		return tokens, errors.NewWrap(errors.ErrInternal, err, "authorize: request build error")
	}

	if err := meta.ToHttpHeaders(req.Header); err != nil {
		return tokens, errors.NewWrap(errors.ErrInvalidArgument, err, "authorize: meta set to request headers error")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return tokens, errors.NewWrap(errors.ErrCanceled, err, "authorize: request canceled error")	
		}
		return tokens, errors.NewWrap(errors.ErrBadGateway, err, "authorize: do request error")
	}
	defer func () { 
		if err = res.Body.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return tokens, errors.NewWrap(errors.ErrInternal, err, "authorize: read response error")
	}

	if res.StatusCode != http.StatusOK {
		kind := mapStatusCodeToKind(res.StatusCode)
		var errResp struct {
			Errors []string `json:"errors"`
		}
		_ = json.Unmarshal(resBody, &errResp)
		msg := "authorize failed"
		if len(errResp.Errors) > 0 {
			msg = errResp.Errors[0]
		}
		return tokens, errors.Newf(kind, "%s (status %d)", msg, res.StatusCode)
	}

	if tokens, err = dto.TokensFromHttpHeaders(res.Header); err != nil {
		return tokens, errors.NewWrapf(errors.ErrUnauthenticated, err, "authorize: invalid tokens in headers")
	}

	return tokens, nil
}
