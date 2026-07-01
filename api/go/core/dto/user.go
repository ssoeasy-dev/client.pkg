package dto

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type UserCompany struct {
	ID                   uuid.UUID
	Name                 string
	ServiceName          *string
	SubscriptionID       *uuid.UUID
	SubscriptionIsActive *bool
}

type User struct {
	Login     string
	Firstname string
	Lastname  string
	Companies []UserCompany
}

func ParseUser(r *http.Response) (User, error) {
	var code User

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&code)
	if err != nil {
		return code, err
	}

	return code, nil
}
