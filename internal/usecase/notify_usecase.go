package usecase

import (
	"context"
	"encoding/json"
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

const flowSendNotification = "SendNotification"

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
	// 1. Validate request
	if err := validator.ValidateNotifyRequest(req); err != nil {
		logger.Error(flowSendNotification, "request validation failed", err,
			"email", req.Email,
			"type", req.Type,
		)
		return nil, err
	}

	// 2. Log incoming request
	logger.Request(flowSendNotification, req)

	// 3. Convert proto request to entity
	history := toEmailHistoryEntity(req)

	// 4. Set initial status
	history.Status = "PENDING"
	history.CreatedAt = time.Now()

	// 5. Save to Repository
	if err := u.repo.SaveEmailHistory(ctx, history); err != nil {
		logger.Error(flowSendNotification, "saving to repository", err,
			"email", history.Email,
			"type", history.Type,
		)
		return nil, fmt.Errorf("failed to save email history: %w", err)
	}
	logger.Info(flowSendNotification, "saved to repository successfully",
		"email", history.Email,
		"type", history.Type,
		"status", history.Status,
	)

	// 6. Marshal and publish to NATS JetStream
	data, err := json.Marshal(history)
	if err != nil {
		logger.Error(flowSendNotification, "marshaling data for NATS", err,
			"email", history.Email,
		)
		return nil, fmt.Errorf("failed to marshal history: %w", err)
	}

	pubAck, err := u.js.Publish(u.subject, data)
	if err != nil {
		logger.Error(flowSendNotification, "publishing to NATS", err,
			"subject", u.subject,
			"email", history.Email,
		)
		return nil, fmt.Errorf("failed to publish to NATS: %w", err)
	}

	// 7. Build and return response
	resp := toNotifyResponse(history)

	logger.Response(flowSendNotification, map[string]any{
		"stream":   pubAck.Stream,
		"sequence": pubAck.Sequence,
		"subject":  u.subject,
		"email":    history.Email,
		"type":     history.Type,
	})

	return resp, nil
}

// toEmailHistoryEntity converts a proto NotifyRequest to an entity.EmailHistory.
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
		t := req.SentAt.AsTime()
		h.SentAt = &t
	}
	return h
}

// toNotifyResponse converts an entity.EmailHistory to a proto NotifyResponse.
func toNotifyResponse(h *entity.EmailHistory) *notifypbv2.NotifyResponse {
	sentAt := timestamppb.New(time.Now())
	if h.SentAt != nil {
		sentAt = timestamppb.New(*h.SentAt)
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
		CreatedAt:      timestamppb.New(h.CreatedAt),
		UpdatedAt:      timestamppb.New(h.UpdatedAt),
	}
}
