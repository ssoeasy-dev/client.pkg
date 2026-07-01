package dto

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ssoeasy-dev/pkg/errors"
)

type Code struct {
	Code      string    `json:"code,omitempty"`
	Verifier  string    `json:"verifier,omitempty"`
	CompanyID uuid.UUID `json:"companyId,omitempty"`
}

var (
	ErrInvalidCode = func(err error) error { return errors.NewWrapf(errors.ErrInvalidArgument, err, "decode code body error") }
	ErrCodeNotProvided = errors.Newf(errors.ErrInvalidArgument, "code not provided")
	ErrVerifierNotProvided = errors.Newf(errors.ErrInvalidArgument, "verifier not provided")
)

func ParseCode(r *http.Request) (Code, error) {
	var code Code

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&code)
	if err != nil {
		return code, ErrInvalidCode(err)
	}

	if code.Code == "" {
		return code, ErrCodeNotProvided
	}

	if code.Verifier == "" {
		return code, ErrVerifierNotProvided
	}

	return code, nil
}
