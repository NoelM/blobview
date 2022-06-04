package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSS3Driver struct {
	client *s3.Client
	config aws.Config
}

func NewAWSS3Driver() *AWSS3Driver {
	return &AWSS3Driver{
		config: aws.Config{},
	}
}

func (a *AWSS3Driver) Start() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	a.config = cfg
	a.client = s3.NewFromConfig(cfg)

	return nil
}

func (a *AWSS3Driver) ListBuckets() (*s3.ListBucketsOutput, error) {
	resp, err := a.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AWSS3Driver) ListObjects(bucket string, prefix string) (*s3.ListObjectsV2Output, error) {
	objectInput := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		MaxKeys:   50,
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}

	list, err := a.client.ListObjectsV2(context.TODO(), objectInput)
	if err != nil {
		return nil, err
	}
	return list, nil
}
