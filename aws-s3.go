package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSS3Driver struct {
	client *s3.S3
}

func NewAWSS3Driver() *AWSS3Driver {
	return &AWSS3Driver{
		client: nil,
	}
}

func (a *AWSS3Driver) Start() error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	a.client = s3.New(sess)

	return nil
}

func (a *AWSS3Driver) ListBuckets() (*s3.ListBucketsOutput, error) {
	resp, err := a.client.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
