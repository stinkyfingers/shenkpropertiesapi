package storage

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	Session *s3.S3
}

const (
	IMAGE_BUCKET = "shenkpropertiesapi"
)

func NewS3(profile string) (*S3, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config: aws.Config{
			Region: aws.String("us-west-1"),
		},
	})
	if err != nil {
		return nil, err
	}

	return &S3{
		Session: s3.New(sess),
	}, nil
}

func (s *S3) List(bucket, prefix string) ([]string, error) {
	var keys []string
	err := s.Session.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, o := range page.Contents {
			keys = append(keys, *o.Key)
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (s *S3) Get(bucket, key string) (io.ReadCloser, error) {
	resp, err := s.Session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				return io.NopCloser(&bytes.Buffer{}), nil
			}
		}
		return nil, err
	}
	return resp.Body, nil
}
