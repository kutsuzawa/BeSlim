package client

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"google.golang.org/api/option"
)

// Database has firebase client.
type Database struct {
	Client *firestore.Client
}

// NewDatabase init Database.
// We get credential.json for firebase from S3, then, we connect firestore.
// Finaly, Database structure obtain *firestore.Client
func NewDatabase() (*Database, error) {
	var credentialPath string
	ctx := context.Background()
	log.Printf("[DEBUG] APP_ENV: %s\n", os.Getenv("APP_ENV"))
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

// getCredentialPathFromS3 download credention file for firebase to 'credentialPath' in lambda.
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
func (db *Database) AddUser(gloupID, userID string) error {
	ctx := context.Background()
	_, _, err := db.Client.Collection("users").Add(ctx, map[string]interface{}{
		"gloupID": gloupID,
		"userID":  userID,
	})
	if err != nil {
		return err
	}
	return nil
}
