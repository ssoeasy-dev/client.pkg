package ssoeasy

import (
	"context"

	"github.com/google/uuid"
)

type Client interface {
	CheckPermission(ctx context.Context, token string, permissionID uuid.UUID) (bool, error)
}
