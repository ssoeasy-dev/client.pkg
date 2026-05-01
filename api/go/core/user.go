package ssoeasy

import (
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id"`
	ServiceID uuid.UUID `json:"service_id"`
	CompanyID uuid.UUID `json:"company_id"`
}

type rawClaims struct {
	UserID    string `json:"user_id"`
	ServiceID string `json:"service_id"`
	CompanyID string `json:"company_id"`
}

func ParseUser(payload []byte) (*User, error) {
	var raw rawClaims
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, NewError("invalid JWT: %w", err)
	}

	userID, err := uuid.Parse(raw.UserID)
	if err != nil {
		return nil, NewError("invalid user_id: %w", err)
	}

	serviceID, err := uuid.Parse(raw.ServiceID)
	if err != nil {
		return nil, NewError("invalid service_id: %w", err)
	}

	companyID, err := uuid.Parse(raw.CompanyID)
	if err != nil {
		return nil, NewError("invalid company_id: %w", err)
	}

	return &User{
		UserID:    userID,
		ServiceID: serviceID,
		CompanyID: companyID,
	}, nil
}
