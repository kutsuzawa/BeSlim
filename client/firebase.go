package client

import (
	"context"
	"errors"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Database has firebase client.
type Database struct {
	Client *firestore.Client
}

// Load has date, weight, and distance.
// Data are obtained from db.
type Load struct {
	Date     time.Time `json:"date" firestore:"date"`
	Weight   float64   `json:"weight" firestore:"weight"`
	Distance float64   `json:"distance" firestore:"distance"`
}

// NewDatabase init Database.
// We get credential.json for firebase from S3, then, we connect firestore.
// Finally, Database structure obtain *firestore.Client
func NewDatabase() (*Database, error) {
	var credentialPath string
	ctx := context.Background()
	if os.Getenv("APP_ENV") != "local" {
		var err error
		credentialPath, err = getCredentialPathFromS3()
		if err != nil {
			return nil, err
		}
	} else {
		credentialPath = "./credential.json"
	}
	opt := option.WithCredentialsFile(credentialPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	db := &Database{Client: client}
	return db, nil
}

// getCredentialPathFromS3 download credential file for firebase to 'credentialPath' in lambda.
func getCredentialPathFromS3() (string, error) {
	// Please put credential.json for firebase to amazon S3.
	// This file is a secret file. So, you set read/write permission for you only.
	// Set env about your bucket name and key (name of credential.json) to serverless.yml.
	bucket := os.Getenv("BUCKET")
	key := os.Getenv("KEY")
	region := os.Getenv("REGION")
	credentialPath := "/tmp/credential.json"

	sess := session.New(&aws.Config{Region: aws.String(region)})
	downloader := s3manager.NewDownloader(sess)
	credentialFile, err := os.Create(credentialPath)
	if err != nil {
		return "", err
	}
	defer credentialFile.Close()
	if _, err := downloader.Download(credentialFile, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return "", err
	}
	return credentialPath, nil
}

// AddLoad add data every user information to firestore.
func (db *Database) AddLoad(userID string, weight, distance float64, date time.Time) error {
	ctx := context.Background()
	load := Load{
		Date:     date,
		Weight:   weight,
		Distance: distance,
	}
	_, _, err := db.Client.Collection("users").Doc(userID).Collection("load").Add(ctx, load)
	if err != nil {
		return err
	}
	return nil
}

// GetDataByUserID execute searching weight and distance data by using userID.
func (db *Database) GetDataByUserID(userID string, start time.Time, end time.Time) ([]Load, error) {
	ctx := context.Background()
	iter := db.Client.Collection("users").Doc(userID).Collection("load").Where("date", ">", start).Where("date", "<", end).OrderBy("date", firestore.Asc).Documents(ctx)
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
