package client

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Database interface {
	AddLoad(userID string, weight, distance float64, date time.Time) error
	GetDataByUserID(userID string, start time.Time, end time.Time) ([]Load, error)
}

// firebase has fb fsClient.
type firebase struct {
	fsClient *firestore.Client
}

func NewFirebase(client *firestore.Client) Database {
	return &firebase{
		fsClient: client,
	}
}

// Load has date, weight, and distance.
// Data are obtained from db.
type Load struct {
	Date     time.Time `json:"date" firestore:"date"`
	Weight   float64   `json:"weight" firestore:"weight"`
	Distance float64   `json:"distance" firestore:"distance"`
}

// AddLoad add data every user information to firestore.
func (db *firebase) AddLoad(userID string, weight, distance float64, date time.Time) error {
	ctx := context.Background()
	load := Load{
		Date:     date,
		Weight:   weight,
		Distance: distance,
	}
	_, _, err := db.fsClient.Collection("users").Doc(userID).Collection("load").Add(ctx, load)
	if err != nil {
		return err
	}
	return nil
}

// GetDataByUserID execute searching weight and distance data by using userID.
func (db *firebase) GetDataByUserID(userID string, start time.Time, end time.Time) ([]Load, error) {
	ctx := context.Background()
	iter := db.fsClient.Collection("users").Doc(userID).Collection("load").Where("date", ">", start).Where("date", "<", end).OrderBy("date", firestore.Asc).Documents(ctx)
	var results []Load
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		result := Load{}
		result.assertion(data)
		results = append(results, result)
	}
	return results, nil
}

func (l *Load) assertion(data map[string]interface{}) (Load, error) {
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
