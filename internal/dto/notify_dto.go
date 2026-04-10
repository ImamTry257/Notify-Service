package dto

import "time"

type NotifyRequestDTO struct {
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Type           string    `json:"type"`
	Data           string    `json:"data"`
	AdditionalData string    `json:"additional_data"`
	Status         string    `json:"status"`
	Metadata       string    `json:"metadata"`
	SentAt         time.Time `json:"sent_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type NotifyResponseDTO struct {
	UUID           string    `json:"uuid"`
	Status         string    `json:"status"`
	Message        string    `json:"message"`
}
