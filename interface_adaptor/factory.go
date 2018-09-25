package interface_adaptor

import (
	"context"

	fb "firebase.google.com/go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/kutsuzawa/slim-load-recorder/driver"
	"google.golang.org/api/option"
)

// Factory is the interface that wraps methods for operating factory of client
type Factory interface {
	Database(config *ClientFactoryConfig) (Database, error)
	Storage(config *ClientFactoryConfig) Storage
}

type factory struct {
	env string
}

// NewClientFactory init factory
func NewClientFactory(env string) Factory {
	return &factory{
		env: env,
	}
}

// ClientFactoryConfig is the config for initializing client
type ClientFactoryConfig struct {
	S3Region string
	S3Bucket string
	S3Key    string
}

// Database init Database interface
func (f *factory) Database(config *ClientFactoryConfig) (Database, error) {
	credentialPath := "./credential.json"
	if f.env != "local" {
		credentialPath = "/tmp/credential.json"
		if err := f.Storage(config).Download(credentialPath); err != nil {
			return nil, err
		}
	}

	ctx := context.Background()
	opt := option.WithCredentialsFile(credentialPath)
	app, err := fb.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	fbClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return driver.NewFirebase(fbClient), nil
}

// Storage init Storage interface
func (f *factory) Storage(config *ClientFactoryConfig) Storage {
	return driver.NewS3(aws.String(config.S3Region), aws.String(config.S3Bucket), aws.String(config.S3Key))
}
