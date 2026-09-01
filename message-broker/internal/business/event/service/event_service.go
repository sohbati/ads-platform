package service

import (
	"context"

	"message-broker/internal/business/event/model"
)

type EventService interface {
	Publish(ctx context.Context, subject string, payload []byte) (*model.PublishResponse, error)
}
