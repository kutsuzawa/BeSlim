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

// Result has date, weight, and distance.
// Data are obtained from db.
type Result struct {
	Date     time.Time `json:"date"`
	Weight   float64   `json:"weight"`
	Distance float64   `json:"distance"`
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

// AddUser add user information to firestore.
func (db *Database) AddUser(userID string, weight, distance float64) error {
	ctx := context.Background()
	_, _, err := db.Client.Collection("users").Doc(userID).Collection("load").Add(ctx, map[string]interface{}{
		"date":     time.Now(),
		"weight":   weight,
		"distance": distance,
	})
	if err != nil {
		return err
	}
	return nil
}

// GetDataByUserID execute searching weight and distance data by using userID.
func (db *Database) GetDataByUserID(userID string) ([]Result, error) {
	ctx := context.Background()
	iter := db.Client.Collection("users").Doc(userID).Collection("load").Where("date", "<", time.Now()).OrderBy("date", firestore.Asc).Documents(ctx)
	var results []Result
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		result := Result{}
		result.assertion(data)
		results = append(results, result)
	}
	return results, nil
}

func (r *Result) assertion(data map[string]interface{}) (Result, error) {
	if date, ok := data["date"].(time.Time); ok {
		r.Date = date
	}
	if r.Date.IsZero() {
		return Result{}, errors.New("failed to assertion")
	}

	if weight, ok := data["weight"].(float64); ok {
		r.Weight = weight
	}
	if distance, ok := data["distance"].(float64); ok {
		r.Distance = distance
	}

	return *r, nil
}
