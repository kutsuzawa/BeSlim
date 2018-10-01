package entity

import (
	"time"

	"github.com/pkg/errors"
)

// Load is data which have relation between User and Result
type Load struct {
	User   User
	Result Result
}

// User has ID
type User struct {
	ID string `json:"user_id" firestore:"user_id"`
}

// Result has date, weight, and distance.
// Data are obtained from db.
type Result struct {
	Date     time.Time `json:"date" firestore:"date"`
	Weight   float64   `json:"weight" firestore:"weight"`
	Distance float64   `json:"distance" firestore:"distance"`
}

// Assertion asserts data
func (l *Result) Assertion(data map[string]interface{}) (Result, error) {
	if date, ok := data["date"].(time.Time); ok {
		l.Date = date
	}
	if l.Date.IsZero() {
		return Result{}, errors.New("failed to assert")
	}

	if weight, ok := data["weight"].(float64); ok {
		l.Weight = weight
	}
	if distance, ok := data["distance"].(float64); ok {
		l.Distance = distance
	}

	return *l, nil
}
