package grpc

import (
	"context"

	"github.com/ImamTry257/lms-proto-notify/gen/go/notify"
	"github.com/ImamTry257/Notify-Service/internal/entity"
	"github.com/ImamTry257/Notify-Service/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotifyHandler struct {
	notify.UnimplementedNotifyServiceServer
	usecase usecase.NotifyUsecase
}

func NewNotifyHandler(usecase usecase.NotifyUsecase) *NotifyHandler {
	return &NotifyHandler{
		usecase: usecase,
	}
}

func (h *NotifyHandler) Send(ctx context.Context, req *notify.NotifyRequest) (*notify.NotifyResponse, error) {
	// Map proto request to entity
	history := &entity.EmailHistory{
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
		history.SentAt = &t
	}

	// Call usecase
	err := h.usecase.SendNotification(ctx, history)
	if err != nil {
		return nil, err
	}

	// Map back to response
	return &notify.NotifyResponse{
		Uuid:           history.IDStr(), // Assuming we want a string version of ID or a real UUID
		Email:          history.Email,
		Phone:          history.Phone,
		Type:           history.Type,
		Data:           history.Data,
		AdditionalData: history.AdditionalData,
		Status:         history.Status,
		SentAt:         timestamppb.New(timeNowOrPtr(history.SentAt)),
		CreatedAt:      timestamppb.New(history.CreatedAt),
		Metadata:       history.Metadata,
		UpdatedAt:      timestamppb.New(history.UpdatedAt),
	}, nil
}

func timeNowOrPtr(t *time.Time) time.Time {
	if t == nil {
		return time.Now()
	}
	return *t
}
