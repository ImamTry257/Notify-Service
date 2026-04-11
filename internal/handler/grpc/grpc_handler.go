package grpc

import (
	"context"

	"github.com/ImamTry257/Notify-Service/internal/usecase"
	notifypbv2 "github.com/ImamTry257/lms-proto-notify/gen/go/notify"
)

type NotifyHandler struct {
	notifypbv2.UnimplementedNotifyServiceServer
	usecase usecase.NotifyUsecase
}

func NewNotifyHandler(usecase usecase.NotifyUsecase) *NotifyHandler {
	return &NotifyHandler{
		usecase: usecase,
	}
}

func (h *NotifyHandler) Send(ctx context.Context, req *notifypbv2.NotifyRequest) (*notifypbv2.NotifyResponse, error) {
	return h.usecase.SendNotification(ctx, req)
}
