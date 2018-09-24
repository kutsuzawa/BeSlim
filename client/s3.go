package client

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3lib "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 has information for connection AWS S3
type S3 struct {
	Region *string
	Bucket *string
	Key    *string
}

// NewS3 init S3
func NewS3(region, bucket, key string) *S3 {
	s3 := &S3{
		Region: aws.String(region),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	return s3
}

// Download execute download file on S3 to lambda.
func (s3 *S3) Download(lambdaPath string) error {
	sess, err := session.NewSession(&aws.Config{Region: s3.Region})
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
		Bucket: s3.Bucket,
		Key:    s3.Key,
	}); err != nil {
		return err
	}
	return nil
}
