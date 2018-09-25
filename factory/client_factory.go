package factory

import (
	"context"

	fb "firebase.google.com/go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/kutsuzawa/slim-load-recorder/client"
	"google.golang.org/api/option"
)

type ClientFactory interface {
	Database(config *ClientFactoryConfig) (client.Database, error)
	Storage(config *ClientFactoryConfig) client.Storage
}

type clientFactory struct {
	env string
}

func NewClientFactory(env string) ClientFactory {
	return &clientFactory{
		env: env,
	}
}

type ClientFactoryConfig struct {
	S3Region string
	S3Bucket string
	S3Key    string
}

// Database init Database interface
func (f *clientFactory) Database(config *ClientFactoryConfig) (client.Database, error) {
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
	return client.NewFirebase(fbClient), nil
}

// Storage init Storage interface
func (*clientFactory) Storage(config *ClientFactoryConfig) client.Storage {
	return client.NewS3(aws.String(config.S3Region), aws.String(config.S3Bucket), aws.String(config.S3Key))
}
