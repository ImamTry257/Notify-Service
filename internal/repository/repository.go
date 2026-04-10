package repository

import (
	"context"
	"github.com/ImamTry257/Notify-Service/internal/entity"
)

type NotifyRepository interface {
	SaveEmailHistory(ctx context.Context, history *entity.EmailHistory) error
	UpdateStatus(ctx context.Context, id int64, status string) error
	GetByID(ctx context.Context, id int64) (*entity.EmailHistory, error)
}
