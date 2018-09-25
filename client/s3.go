package client

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3lib "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage interface {
	Download(lambdaPath string) error
}

// s3 has information for connection AWS s3
type s3 struct {
	region *string
	bucket *string
	key    *string
}

func NewS3(region, bucket, key *string) Storage {
	return &s3{
		region: region,
		bucket: bucket,
		key:    key,
	}
}

// Download execute download file on s3 to lambda.
func (s3 *s3) Download(lambdaPath string) error {
	sess, err := session.NewSession(&aws.Config{Region: s3.region})
	if err != nil {
		return err
	}
	downloader := s3manager.NewDownloader(sess)
	file, err := os.Create(lambdaPath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := downloader.Download(file, &s3lib.GetObjectInput{
		Bucket: s3.bucket,
		Key:    s3.key,
	}); err != nil {
		return err
	}
	return nil
}
