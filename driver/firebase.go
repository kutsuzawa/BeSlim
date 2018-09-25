package driver

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kutsuzawa/slim-load-recorder/application"
	"google.golang.org/api/iterator"
)

// FireBase is the interface wrap methods for operating firebase
type FireBase interface {
	AddLoad(userID string, weight, distance float64, date time.Time) error
	GetDataByUserID(userID string, start time.Time, end time.Time) ([]application.Load, error)
}

// firebase has fb fsClient.
type firebase struct {
	fsClient *firestore.Client
}

// NewFirebase init firebase
func NewFirebase(client *firestore.Client) FireBase {
	return &firebase{
		fsClient: client,
	}
}

// AddLoad add data every user information to firestore.
func (db *firebase) AddLoad(userID string, weight, distance float64, date time.Time) error {
	ctx := context.Background()
	load := application.Load{
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
func (db *firebase) GetDataByUserID(userID string, start time.Time, end time.Time) ([]application.Load, error) {
	ctx := context.Background()
	iter := db.fsClient.Collection("users").Doc(userID).Collection("load").Where("date", ">", start).Where("date", "<", end).OrderBy("date", firestore.Asc).Documents(ctx)
	var results []application.Load
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		result := application.Load{}
		result.Assertion(data)
		results = append(results, result)
	}
	return results, nil
}
