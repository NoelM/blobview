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

func (a *AWSS3Driver) listBucketsInternal() (*s3.ListBucketsOutput, error) {
	resp, err := a.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AWSS3Driver) ListBuckets() (*ObjectList, error) {
	raw, err := a.listBucketsInternal()
	if err != nil {
		return nil, err
	}

	objectList := convertListBucketOutput(raw)
	return objectList, nil
}

func (a *AWSS3Driver) listObjectsInternal(bucket, prefix string) (*s3.ListObjectsV2Output, error) {
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

func (a *AWSS3Driver) ListObjects(bucket, prefix string) (*ObjectList, error) {
	raw, err := a.listObjectsInternal(bucket, prefix)
	if err != nil {
		return nil, err
	}

	objectList := convertListObjectOutput(raw)
	objectList.Bucket = bucket
	objectList.Prefix = prefix
	objectList.Token = *raw.ContinuationToken

	return objectList, nil
}

func (a *AWSS3Driver) listObjectsNextInternal(bucket, prefix, token string) (*s3.ListObjectsV2Output, error) {
	objectInput := &s3.ListObjectsV2Input{
		Bucket:            aws.String(bucket),
		MaxKeys:           50,
		Prefix:            aws.String(prefix),
		Delimiter:         aws.String("/"),
		ContinuationToken: aws.String(token),
	}

	list, err := a.client.ListObjectsV2(context.TODO(), objectInput)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a *AWSS3Driver) ListObjectsNext(bucket, prefix, token string) (*ObjectList, error) {
	raw, err := a.listObjectsNextInternal(bucket, prefix, token)
	if err != nil {
		return nil, err
	}

	objectList := convertListObjectOutput(raw)
	objectList.Bucket = bucket
	objectList.Prefix = prefix
	objectList.Token = *raw.ContinuationToken

	return objectList, nil
}

func convertListBucketOutput(lb *s3.ListBucketsOutput) *ObjectList {
	var objectList *ObjectList

	for _, bucket := range lb.Buckets {
		objectList.List = append(objectList.List, Object{
			Key:  *bucket.Name,
			Date: *bucket.CreationDate,
			Type: BucketType,
		})
	}

	return objectList
}

func convertListObjectOutput(lo *s3.ListObjectsV2Output) *ObjectList {
	var objectList *ObjectList

	for _, dir := range lo.CommonPrefixes {
		objectList.List = append(objectList.List, Object{
			Key:  *dir.Prefix,
			Type: DirectoryType,
		})
	}
	for _, obj := range lo.Contents {
		objectList.List = append(objectList.List, Object{
			Key:   *obj.Key,
			Date:  *obj.LastModified,
			Bytes: obj.Size,
			Type:  FileType,
		})
	}

	return objectList
}
