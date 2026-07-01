package dto

import (
	"net/http"

	"github.com/ssoeasy-dev/pkg/errors"
)

const (
	AccessHeader  = "Authorization"
	RefreshHeader = "Refresh"
)

type Tokens struct {
	Access  string `json:"access,omitempty"`
	Refresh string `json:"refresh,omitempty"`
	Payload Payload
}

func (t Tokens) ToHttpHeaders(h http.Header) error {
	if t.Access == "" {
		return errors.New(errors.ErrUnauthenticated, "access token is empty")
	}
	h.Set(AccessHeader, BearerPrefix+t.Access)

	if t.Refresh == "" {
		return errors.New(errors.ErrUnauthenticated, "refresh token is empty")
	}
	h.Set(RefreshHeader, t.Refresh)

	return nil
}

func TokensFromHttpHeaders(h http.Header) (Tokens, error) {
	var tokens Tokens

	tokens.Access = h.Get(AccessHeader)
	if tokens.Access == "" {
		return tokens, errors.New(errors.ErrUnauthenticated, "access token header missing")
	}
	tokens.Access = ExtractBearer(tokens.Access)

	tokens.Refresh = h.Get(RefreshHeader)
	if tokens.Refresh == "" {
		return tokens, errors.New(errors.ErrUnauthenticated, "refresh token header missing")
	}

    payload, err := PayloadFromTokens(tokens)
    if err != nil {
        return tokens, err
    }
    tokens.Payload = payload

	return tokens, nil
}
