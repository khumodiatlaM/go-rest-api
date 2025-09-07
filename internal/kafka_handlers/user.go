package kafka_handlers

import (
	"context"
	"encoding/json"
	"go-rest-api/internal/core"
	"go-rest-api/pkg/kafka"
	"go-rest-api/pkg/logger"
)

type UserEventService struct {
	producer *kafka.Producer
	logger   logger.CustomLogger
	topic    string
}

func NewUserEventService(producer *kafka.Producer, logger logger.CustomLogger, topic string) core.UserEventService {
	return UserEventService{
		producer: producer,
		logger:   logger,
		topic:    topic,
	}
}

func (s UserEventService) PublishUserCreatedEvent(ctx context.Context, user *core.User) error {
	userDataBytes, err := json.Marshal(user)
	if err != nil {
		s.logger.Error("failed to marshal user data: ", err)
		return err
	}
	err = s.producer.Produce(s.topic, user.ID.String(), userDataBytes)
	if err != nil {
		s.logger.Error("failed to produce user created event: ", err)
		return err
	}
	s.logger.Info("user created event successfully published to topic: ", s.topic)
	return nil
}
