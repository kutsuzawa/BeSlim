package usecase

import (
	"time"
)

// Request is used to receive requests
type Request struct {
	UserID   string  `json:"user_id"`
	Weight   float64 `json:"weight"`
	Distance float64 `json:"distance"`
	Date     string  `json:"date"`
	StartAt  string  `json:"start_at"`
	EndAt    string  `json:"end_at"`
}

func (r *Request) date() (time.Time, error) {
	return r.parseStrToTime(r.Date)
}

func (r *Request) startAt() (time.Time, error) {
	return r.parseStrToTime(r.StartAt)
}

func (r *Request) endAt() (time.Time, error) {
	return r.parseStrToTime(r.EndAt)
}

func (r *Request) parseStrToTime(str string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
