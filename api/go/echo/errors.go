package echo

import (
    "net/http"
    "github.com/ssoeasy-dev/pkg/errors"
)

func httpStatusFromError(err error) int {
    kind := errors.Kind(err)
    switch kind {
    case errors.ErrInvalidArgument, errors.ErrFailedPrecondition,
        errors.ErrUnprocessableEntity, errors.ErrMethodNotAllowed,
        errors.ErrNotAcceptable, errors.ErrUnsupportedMediaType,
        errors.ErrRangeNotSatisfiable, errors.ErrExpectationFailed,
        errors.ErrPreconditionRequired, errors.ErrRequestHeaderFieldsTooLarge:
        return http.StatusBadRequest

    case errors.ErrUnauthenticated:
        return http.StatusUnauthorized

    case errors.ErrPermissionDenied:
        return http.StatusForbidden

    case errors.ErrNotFound, errors.ErrGone:
        return http.StatusNotFound

    case errors.ErrAlreadyExists, errors.ErrConflict:
        return http.StatusConflict

    case errors.ErrTooManyRequests:
        return http.StatusTooManyRequests

    case errors.ErrCanceled, errors.ErrRequestTimeout:
        return http.StatusRequestTimeout

    case errors.ErrBadGateway:
        return http.StatusBadGateway

    case errors.ErrUnavailable, errors.ErrGatewayTimeout:
        return http.StatusServiceUnavailable

    case errors.ErrInternal, errors.ErrDataLoss, errors.ErrUnknown:
        return http.StatusInternalServerError

    default:
        return http.StatusInternalServerError
    }
}
