package client

import (
	"net/http"

	"github.com/ssoeasy-dev/pkg/errors"
)

// mapStatusCodeToKind преобразует HTTP-статус в соответствующий вид ошибки.
func mapStatusCodeToKind(status int) error {
	switch status {
	case http.StatusBadRequest:
		return errors.ErrInvalidArgument
	case http.StatusUnauthorized:
		return errors.ErrUnauthenticated
	case http.StatusForbidden:
		return errors.ErrPermissionDenied
	case http.StatusNotFound:
		return errors.ErrNotFound
	case http.StatusConflict:
		return errors.ErrConflict
	case http.StatusGone:
		return errors.ErrGone
	case http.StatusTooManyRequests:
		return errors.ErrTooManyRequests
	case http.StatusInternalServerError:
		return errors.ErrInternal
	case http.StatusBadGateway:
		return errors.ErrBadGateway
	case http.StatusServiceUnavailable:
		return errors.ErrUnavailable
	case http.StatusGatewayTimeout:
		return errors.ErrGatewayTimeout
	default:
		return errors.ErrUnknown
	}
}
