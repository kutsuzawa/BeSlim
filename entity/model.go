package entity

import (
	"time"

	"github.com/pkg/errors"
)

// Load has date, weight, and distance.
// Data are obtained from db.
type Load struct {
	Date     time.Time `json:"date" firestore:"date"`
	Weight   float64   `json:"weight" firestore:"weight"`
	Distance float64   `json:"distance" firestore:"distance"`
}

// Assertion asserts data
func (l *Load) Assertion(data map[string]interface{}) (Load, error) {
	if date, ok := data["date"].(time.Time); ok {
		l.Date = date
	}
	if l.Date.IsZero() {
		return Load{}, errors.New("failed to assert")
	}

	if weight, ok := data["weight"].(float64); ok {
		l.Weight = weight
	}
	if distance, ok := data["distance"].(float64); ok {
		l.Distance = distance
	}

	return *l, nil
}
