package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ImamTry257/Notify-Service/internal/entity"
)

type mysqlNotifyRepository struct {
	db *sql.DB
}

func NewMySQLNotifyRepository(db *sql.DB) NotifyRepository {
	return &mysqlNotifyRepository{db: db}
}

func (r *mysqlNotifyRepository) SaveEmailHistory(ctx context.Context, history *entity.EmailHistory) error {
	query := `INSERT INTO email_histories (email, phone, type, data, additional_data, status, metadata, sent_at, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	res, err := r.db.ExecContext(ctx, query,
		history.Email,
		history.Phone,
		history.Type,
		history.Data,
		history.AdditionalData,
		history.Status,
		history.Metadata,
		history.SentAt,
		now,
		now,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	history.ID = id
	history.CreatedAt = now
	history.UpdatedAt = now
	return nil
}

func (r *mysqlNotifyRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE email_histories SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *mysqlNotifyRepository) GetByID(ctx context.Context, id int64) (*entity.EmailHistory, error) {
	query := `SELECT id, email, phone, type, data, additional_data, status, metadata, sent_at, created_at, updated_at 
			  FROM email_histories WHERE id = ?`
	
	var h entity.EmailHistory
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&h.ID, &h.Email, &h.Phone, &h.Type, &h.Data, &h.AdditionalData, &h.Status, &h.Metadata, &h.SentAt, &h.CreatedAt, &h.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
