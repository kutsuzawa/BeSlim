package firebase

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/kutsuzawa/slim-load-recorder/entity"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Firebase has client for connecting firebase
type Firebase struct {
	client *firestore.Client
}

// NewFirebase init Firebase struct
func NewFirebase(env, region, bucket, key string) (*Firebase, error) {
	fb := &Firebase{}
	s3 := NewS3(&region, &bucket, &key)
	credentialPath := "./credential.json"
	if env != "local" {
		credentialPath = "/tmp/credential.json"
		if err := s3.Download(credentialPath); err != nil {
			return nil, err
		}
	}

	ctx := context.Background()
	opt := option.WithCredentialsFile(credentialPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	fb.client, err = app.Firestore(ctx)
	return fb, nil
}

// Add is implementation of DatabaseDriver in adapter package.
// It add Load every user to Firebase.
func (db *Firebase) Add(userID string, load entity.Load) error {
	ctx := context.Background()
	_, _, err := db.client.Collection("users").Doc(userID).Collection("load").Add(ctx, load)
	if err != nil {
		return err
	}
	return nil
}

// Search is implementation of DatabaseDriver in adapter package.
// It search Loads that satisfy desired duration from Firebase.
func (db *Firebase) Search(userID string, start, end time.Time) ([]entity.Load, error) {
	ctx := context.Background()
	iter := db.client.Collection("users").Doc(userID).Collection("load").Where("date", ">", start).Where("date", "<", end).OrderBy("date", firestore.Asc).Documents(ctx)
	var results []entity.Load
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		result := entity.Load{}
		// TODO: AssertionはAdapterレイヤーでやる気がする.
		result.Assertion(data)
		results = append(results, result)
	}
	return results, nil

}
