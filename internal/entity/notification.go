package entity

import (
	"fmt"
	"time"
)

const (
	TYPE_OTP            = "OTP"
	TYPE_RESET_PASSWORD = "RESET_PASSWORD"
	TYPE_ACTIVATION     = "ACTIVATION"
)

type EmailHistory struct {
	ID             int64     `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	Phone          string    `json:"phone" db:"phone"`
	Type           string    `json:"type" db:"type"`
	Data           string    `json:"data" db:"data"`
	AdditionalData string    `json:"additional_data" db:"additional_data"`
	Status         string    `json:"status" db:"status"`
	Metadata       string    `json:"metadata" db:"metadata"`
	SentAt         *time.Time `json:"sent_at" db:"sent_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func (e *EmailHistory) IDStr() string {
	return fmt.Sprintf("%d", e.ID)
}
