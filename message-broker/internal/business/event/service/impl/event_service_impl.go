package impl

import (
	"context"

	"message-broker/internal/business/event/errorcode"
	"message-broker/internal/business/event/model"
	"message-broker/internal/business/event/service"
	"message-broker/internal/core/broker"
	"message-broker/internal/core/exception"
)

type eventService struct {
	broker *broker.Connection
}

func NewEventService(conn *broker.Connection) service.EventService {
	return &eventService{broker: conn}
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

	if !s.broker.IsConnected() {
		return nil, exception.NewAppError(
			errorcode.ErrBrokerUnavailable.Code, errorcode.ErrBrokerUnavailable.HttpStatus, subject)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := s.broker.Publish(subject, payload); err != nil {
		return nil, exception.NewAppError(
			errorcode.ErrPublishFailed.Code, errorcode.ErrPublishFailed.HttpStatus, subject).WithCause(err)
	}

	return &model.PublishResponse{
		Subject: subject,
		Message: "published",
	}, nil
}
