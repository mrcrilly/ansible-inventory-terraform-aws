package main

import (
	"errors"
	"os"
)

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func downloadStateFromS3(bucket, key, toFile string) error {
	s3Region := os.Getenv("AWS_DEFAULT_REGION")

	if s3Region == "" {
		return errors.New("AWS region is required in environment variable: AWS_DEFAULT_REGION.")
	}

	file, err := os.Create(toFile)

	if err != nil {
		return err
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(s3Region)}))

	_o := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err = downloader.Download(file, _o)

	if err != nil {
		return err
	}

	return nil
}
