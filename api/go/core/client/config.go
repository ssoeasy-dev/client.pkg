package client

import (
	"time"

	"github.com/google/uuid"
)

type Config struct {
	BaseURL   string
	ServiceID uuid.UUID
	Timeout   time.Duration
}
