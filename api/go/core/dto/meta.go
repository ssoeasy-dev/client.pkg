package dto

import (
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/ssoeasy-dev/pkg/errors"
)

const ( 
	TraceIDHeader = "x-trace-id"
	OriginHeader = "origin"
)

// Метаданные запросов. Включает общие заголовки и прочие метаданные.
// - Origin (url.URL): Домен клиента
// - TraceID (uuid.UUID): Идентификатор для отслеживания запросов
type Meta struct {
    Origin  *url.URL
    TraceID uuid.UUID
}

func (m Meta) ToHttpHeaders(h http.Header) error {
    h.Set(TraceIDHeader, m.TraceID.String())
    if m.Origin != nil {
        h.Set(OriginHeader, m.Origin.String())
    }
    h.Set("Content-Type", "application/json")
    return nil
}

func MetaFromHttpHeaders(h http.Header) (Meta, error) {
    var meta Meta

    if traceID := h.Get(TraceIDHeader); traceID != "" {
        id, err := uuid.Parse(traceID)
        if err != nil {
            return meta, errors.NewWrap(errors.ErrInvalidArgument, err, "invalid trace-id")
        }
        meta.TraceID = id
    }

    if origin := h.Get(OriginHeader); origin != "" {
        u, err := url.Parse(origin)
        if err != nil {
            return meta, errors.NewWrap(errors.ErrInvalidArgument, err, "invalid origin")
        }
        meta.Origin = u
    }

    return meta, nil
}
