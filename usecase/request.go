package usecase

import (
	"time"
)

// Request is used to receive requests
type Request struct {
	UserID   string    `json:"user_id"`
	Weight   float64   `json:"weight"`
	Distance float64   `json:"distance"`
	Date     time.Time `json:"date"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}
