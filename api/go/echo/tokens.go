package echo

import (
	"net/http"

	"github.com/ssoeasy-dev/client.pkg/api/go/core/v2/dto"
	"github.com/ssoeasy-dev/pkg/errors"
)

type Tokens struct {
	dto.Tokens
}

func TokensFromCore(tokens dto.Tokens) Tokens {
	return Tokens{tokens}
}

func (t Tokens) ToHttpCookie(w http.ResponseWriter, cfg CookieConfig) error {
	if t.Access == "" {
		return errors.New(errors.ErrUnauthenticated, "access token is empty")
	}
    w.Header().Set(dto.AccessHeader, dto.BearerPrefix+t.Access)

	if t.Refresh == "" {
        return nil
    }

    http.SetCookie(w, &http.Cookie{
        Name:     dto.RefreshHeader,
        Value:    t.Refresh,
        Path:     cfg.Path,
        Domain:   cfg.Domain,
        MaxAge:   cfg.MaxAge,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    })

	return nil
}

func TokensFromHttpCookie(r *http.Request) (dto.Tokens, error) {
	var tokens dto.Tokens

	tokens.Access = r.Header.Get(dto.AccessHeader)
	if tokens.Access == "" {
		return tokens, errors.New(errors.ErrUnauthenticated, "access token header missing")
	}
	tokens.Access = dto.ExtractBearer(tokens.Access)

	cookie, err := r.Cookie(dto.RefreshHeader)
	if err != nil {
		return tokens, errors.NewWrap(errors.ErrUnauthenticated, err, "refresh cookie missing")
	}
	tokens.Refresh = cookie.Value

    payload, err := dto.PayloadFromTokens(tokens)
    if err != nil {
        return tokens, err
    }
    tokens.Payload = payload

	return tokens, nil
}
