package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ImamTry257/Notify-Service/internal/entity"
	"github.com/ImamTry257/Notify-Service/internal/repository"
	"github.com/ImamTry257/Notify-Service/pkg/logger"
	"github.com/nats-io/nats.go"
)

const flowSendNotification = "SendNotification"

type NotifyUsecase interface {
	SendNotification(ctx context.Context, request *entity.EmailHistory) error
}

type notifyUsecase struct {
	repo    repository.NotifyRepository
	js      nats.JetStreamContext
	subject string
}

func NewNotifyUsecase(repo repository.NotifyRepository, js nats.JetStreamContext, subject string) NotifyUsecase {
	return &notifyUsecase{
		repo:    repo,
		js:      js,
		subject: subject,
	}
}

func (u *notifyUsecase) SendNotification(ctx context.Context, history *entity.EmailHistory) error {
	// 1. Log incoming request
	logger.Request(flowSendNotification, history)

	// 2. Initial status
	history.Status = "PENDING"
	history.CreatedAt = time.Now()

	// 3. Save to Repository
	if err := u.repo.SaveEmailHistory(ctx, history); err != nil {
		logger.Error(flowSendNotification, "saving to repository", err,
			"email", history.Email,
			"type", history.Type,
		)
		return fmt.Errorf("failed to save email history: %w", err)
	}
	logger.Info(flowSendNotification, "saved to repository successfully",
		"email", history.Email,
		"type", history.Type,
		"status", history.Status,
	)

	// 4. Marshal and publish to NATS JetStream
	data, err := json.Marshal(history)
	if err != nil {
		logger.Error(flowSendNotification, "marshaling data for NATS", err,
			"email", history.Email,
		)
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	pubAck, err := u.js.Publish(u.subject, data)
	if err != nil {
		logger.Error(flowSendNotification, "publishing to NATS", err,
			"subject", u.subject,
			"email", history.Email,
		)
		return fmt.Errorf("failed to publish to NATS: %w", err)
	}

	// 5. Log success response
	logger.Response(flowSendNotification, map[string]any{
		"stream":   pubAck.Stream,
		"sequence": pubAck.Sequence,
		"subject":  u.subject,
		"email":    history.Email,
		"type":     history.Type,
	})

	return nil
}
