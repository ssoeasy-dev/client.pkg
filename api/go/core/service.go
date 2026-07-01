package core

import (
	"context"

	"github.com/google/uuid"
	"github.com/ssoeasy-dev/client.pkg/api/go/core/v2/client"
	"github.com/ssoeasy-dev/client.pkg/api/go/core/v2/dto"
)

type Service struct {
	client *client.Client
}

func NewService(cfg client.Config) (*Service, error) {
	client, err := client.NewClient(cfg);
	if err != nil {
		return nil, err
	}

	return &Service{
		client: client,
	}, nil
}

type Callback = func(ctx context.Context, user dto.Payload) error

func (s *Service) Authorize(ctx context.Context, meta dto.Meta, code dto.Code, cb Callback) (dto.Tokens, error) {
	tokens, err := s.client.Authorize(ctx, meta, code)
	if err != nil {
		return tokens, err
	}

	if cb != nil {
		err := cb(ctx, tokens.Payload)
		return tokens, err
	}

	return tokens, nil
}

func (s *Service) Check(ctx context.Context, meta dto.Meta, tokens dto.Tokens, permissionID uuid.UUID, cb Callback) error {
	err := s.client.Check(ctx, meta, tokens, permissionID)
	if err != nil {
		return err
	}

	if cb != nil {
		if err := cb(ctx, tokens.Payload); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Logout(ctx context.Context, meta dto.Meta, tokens dto.Tokens, cb Callback) error {
	err := s.client.Logout(ctx, meta, tokens)
	if err != nil {
		return err
	}

	if cb != nil {
		if err := cb(ctx, tokens.Payload); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Me(ctx context.Context, meta dto.Meta, tokens dto.Tokens, cb Callback) (dto.User, error) {
	user, err := s.client.Me(ctx, meta, tokens)
	if err != nil {
		return user, err
	}

	if cb != nil {
		if err := cb(ctx, tokens.Payload); err != nil {
			return user, err
		}
	}

	return user, err
}

func (s *Service) Refresh(ctx context.Context, meta dto.Meta, tokens dto.Tokens, cb Callback) (dto.Tokens, error) {
	tokens, err := s.client.Refresh(ctx, meta, tokens)
	if err != nil {
		return tokens, err
	}

	if cb != nil {
		if err := cb(ctx, tokens.Payload); err != nil {
			return tokens, err
		}
	}

	return tokens, err
}
