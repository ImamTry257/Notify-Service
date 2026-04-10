package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/ImamTry257/Notify-Service/internal/entity"
	"github.com/ImamTry257/Notify-Service/internal/repository"
)

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
	// 1. Initial status
	history.Status = "PENDING"
	history.CreatedAt = time.Now()

	// 2. Save to Repository
	if err := u.repo.SaveEmailHistory(ctx, history); err != nil {
		return err
	}

	// 3. Publish to NATS Jetstream
	data, err := json.Marshal(history)
	if err != nil {
		return err
	}

	_, err = u.js.Publish(u.subject, data)
	return err
}
