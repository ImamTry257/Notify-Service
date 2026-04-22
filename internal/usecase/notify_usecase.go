```go
package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ImamTry257/Notify-Service/internal/entity"
	"github.com/ImamTry257/Notify-Service/internal/repository"
	"github.com/ImamTry257/Notify-Service/pkg/logger"
	"github.com/ImamTry257/Notify-Service/pkg/validator"
	notifypbv2 "github.com/ImamTry257/lms-proto-notify/gen/go/notify"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	flowSendNotification = "SendNotification"
	statusPending        = "PENDING"
)

var (
	ErrNilNotifyRequest = errors.New("notify request is nil")
	ErrEmptySubject     = errors.New("nats subject is empty")
)

type NotifyUsecase interface {
	SendNotification(ctx context.Context, req *notifypbv2.NotifyRequest) (*notifypbv2.NotifyResponse, error)
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

func (u *notifyUsecase) SendNotification(ctx context.Context, req *notifypbv2.NotifyRequest) (*notifypbv2.NotifyResponse, error) {
	if req == nil {
		logger.Error(flowSendNotification, "request is nil", ErrNilNotifyRequest)
		return nil, ErrNilNotifyRequest
	}

	if u.subject == "" {
		logger.Error(flowSendNotification, "nats subject is empty", ErrEmptySubject)
		return nil, ErrEmptySubject
	}

	if err := validator.ValidateNotifyRequest(req); err != nil {
		logger.Error(flowSendNotification, "request validation failed", err,
			"email", req.Email,
			"type", req.Type,
		)
		return nil, err
	}

	logger.Request(flowSendNotification, req)

	history := toEmailHistoryEntity(req)
	now := time.Now().UTC()
	history.Status = statusPending
	history.CreatedAt = now
	history.UpdatedAt = now

	if err := u.repo.SaveEmailHistory(ctx, history); err != nil {
		logger.Error(flowSendNotification, "saving to repository failed", err,
			"email", history.Email,
			"type", history.Type,
		)
		return nil, fmt.Errorf("save email history: %w", err)
	}

	pubAck, err := u.publishHistory(ctx, history)
	if err != nil {
		return nil, err
	}

	resp := toNotifyResponse(history)

	logger.Response(flowSendNotification, map[string]any{
		"stream":   pubAck.Stream,
		"sequence": pubAck.Sequence,
		"subject":  u.subject,
		"email":    history.Email,
		"type":     history.Type,
		"status":   history.Status,
	})

	return resp, nil
}

func (u *notifyUsecase) publishHistory(ctx context.Context, history *entity.EmailHistory) (*nats.PubAck, error) {
	payload, err := json.Marshal(history)
	if err != nil {
		logger.Error(flowSendNotification, "marshal history failed", err,
			"email", history.Email,
			"type", history.Type,
		)
		return nil, fmt.Errorf("marshal history: %w", err)
	}

	msg := &nats.Msg{
		Subject: u.subject,
		Data:    payload,
	}

	pubAck, err := u.js.PublishMsg(msg, nats.Context(ctx))
	if err != nil {
		logger.Error(flowSendNotification, "publish to nats failed", err,
			"subject", u.subject,
			"email", history.Email,
			"type", history.Type,
		)
		return nil, fmt.Errorf("publish to nats: %w", err)
	}

	if pubAck == nil {
		err = errors.New("nil publish ack")
		logger.Error(flowSendNotification, "publish ack is nil", err,
			"subject", u.subject,
			"email", history.Email,
			"type", history.Type,
		)
		return nil, err
	}

	return pubAck, nil
}

func toEmailHistoryEntity(req *notifypbv2.NotifyRequest) *entity.EmailHistory {
	h := &entity.EmailHistory{
		Email:          req.Email,
		Phone:          req.Phone,
		Type:           req.Type,
		Data:           req.Data,
		AdditionalData: req.AdditionalData,
		Status:         req.Status,
		Metadata:       req.Metadata,
	}

	if req.SentAt != nil {
		t := req.SentAt.AsTime().UTC()
		h.SentAt = &t
	}

	return h
}

func toNotifyResponse(h *entity.EmailHistory) *notifypbv2.NotifyResponse {
	var sentAt *timestamppb.Timestamp
	if h.SentAt != nil {
		sentAt = timestamppb.New(h.SentAt.UTC())
	} else {
		sentAt = timestamppb.New(time.Now().UTC())
	}

	return &notifypbv2.NotifyResponse{
		Uuid:           h.IDStr(),
		Email:          h.Email,
		Phone:          h.Phone,
		Type:           h.Type,
		Data:           h.Data,
		AdditionalData: h.AdditionalData,
		Status:         h.Status,
		Metadata:       h.Metadata,
		SentAt:         sentAt,
		CreatedAt:      toProtoTimestamp(h.CreatedAt),
		UpdatedAt:      toProtoTimestamp(h.UpdatedAt),
	}
}

func toProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t.UTC())
}
```