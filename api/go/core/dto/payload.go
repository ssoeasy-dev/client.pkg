package dto

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/ssoeasy-dev/pkg/errors"
)

type Payload struct {
	UserID    uuid.UUID
	ServiceID uuid.UUID
	CompanyID uuid.UUID
}

const TokenPayloadContextKey = "ssoeasy_payload"

func (p Payload) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, TokenPayloadContextKey, p)
}

func PayloadFromContext(ctx context.Context) (Payload, error) {
	payload, ok := ctx.Value(TokenPayloadContextKey).(Payload)
	if !ok {
		return payload, errors.ErrPayloadTooLarge
	}
	return payload, nil
}

func PayloadFromTokens(tokens Tokens) (Payload, error) {
	var payload Payload

	parts := strings.Split(tokens.Access, ".")
	if len(parts) != 3 {
		return payload, errors.Newf(errors.ErrInvalidArgument, "invalid JWT: expected 3 parts")
	}
	rawByte, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return payload, errors.NewWrapf(errors.ErrInvalidArgument, err, "invalid JWT: base64 decode failed")
	}

	var raw *struct {
		UserID    string `json:"user_id"`
		ServiceID string `json:"service_id"`
		CompanyID string `json:"company_id"`
	}

	if err := json.Unmarshal(rawByte, &raw); err != nil {
		return payload, errors.NewWrapf(errors.ErrInvalidArgument, err, "invalid JWT: JSON unmarshal failed")
	}

	payload.UserID, err = uuid.Parse(raw.UserID)
	if err != nil {
		return payload, errors.NewWrapf(errors.ErrInvalidArgument, err, "invalid user_id in token")
	}

	payload.ServiceID, err = uuid.Parse(raw.ServiceID)
	if err != nil {
		return payload, errors.NewWrapf(errors.ErrInvalidArgument, err, "invalid service_id in token")
	}

	payload.CompanyID, err = uuid.Parse(raw.CompanyID)
	if err != nil {
		return payload, errors.NewWrapf(errors.ErrInvalidArgument, err, "invalid company_id in token")
	}

	return payload, nil
}
