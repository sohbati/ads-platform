package impl

import (
	"context"

	"nats-message-broker/internal/business/event/errorcode"
	"nats-message-broker/internal/business/event/model"
	"nats-message-broker/internal/business/event/service"
	"nats-message-broker/internal/core/exception"
	"nats-message-broker/internal/core/natsconn"
)

type eventService struct {
	nats *natsconn.Connection
}

func NewEventService(nats *natsconn.Connection) service.EventService {
	return &eventService{nats: nats}
}

func (s *eventService) Publish(ctx context.Context, subject string, payload []byte) (*model.PublishResponse, error) {
	if subject == "" {
		return nil, exception.NewAppError(
			errorcode.ErrSubjectEmpty.Code, errorcode.ErrSubjectEmpty.HttpStatus, subject)
	}

	if len(payload) == 0 {
		return nil, exception.NewAppError(
			errorcode.ErrPayloadEmpty.Code, errorcode.ErrPayloadEmpty.HttpStatus, subject)
	}

	if !s.nats.IsConnected() {
		return nil, exception.NewAppError(
			errorcode.ErrNatsUnavailable.Code, errorcode.ErrNatsUnavailable.HttpStatus, subject)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := s.nats.Publish(subject, payload); err != nil {
		return nil, exception.NewAppError(
			errorcode.ErrPublishFailed.Code, errorcode.ErrPublishFailed.HttpStatus, subject).WithCause(err)
	}

	return &model.PublishResponse{
		Subject: subject,
		Message: "published",
	}, nil
}
