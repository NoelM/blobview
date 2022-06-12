package main

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
	"strings"
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

	objectList := convertListObjectOutput(bucket, raw)
	objectList.Bucket = bucket
	objectList.Prefix = prefix
	if raw.ContinuationToken != nil {
		objectList.Token = *raw.ContinuationToken
	}

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

	objectList := convertListObjectOutput(bucket, raw)
	objectList.Bucket = bucket
	objectList.Prefix = prefix
	objectList.Token = *raw.ContinuationToken

	return objectList, nil
}

func convertListBucketOutput(lb *s3.ListBucketsOutput) *ObjectList {
	objectList := &ObjectList{}

	for _, bucket := range lb.Buckets {
		objectList.List = append(objectList.List, Object{
			Provider: AWS,
			Key:      *bucket.Name,
			PrintKey: *bucket.Name,
			Date:     *bucket.CreationDate,
			Type:     BucketType,
		})
	}

	return objectList
}

func convertListObjectOutput(bucket string, lo *s3.ListObjectsV2Output) *ObjectList {
	objectList := &ObjectList{}

	for _, dir := range lo.CommonPrefixes {
		prefixList := strings.Split(*dir.Prefix, "/")
		var printKey string
		if len(prefixList) >= 2 {
			printKey = prefixList[len(prefixList)-2]
		}

		objectList.List = append(objectList.List, Object{
			Provider: AWS,
			Bucket:   bucket,
			Key:      *dir.Prefix,
			PrintKey: printKey,
			Type:     DirectoryType,
		})
	}
	for _, obj := range lo.Contents {
		prefixList := strings.Split(*obj.Key, "/")
		var printKey string
		if len(prefixList) >= 1 {
			printKey = prefixList[len(prefixList)-1]
		}

		objectList.List = append(objectList.List, Object{
			Provider: AWS,
			Bucket:   bucket,
			Key:      *obj.Key,
			PrintKey: printKey,
			Date:     *obj.LastModified,
			Bytes:    obj.Size,
			Type:     FileType,
		})
	}

	return objectList
}

func (a *AWSS3Driver) Download(object Object, destination string) error {
	output, err := a.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	})
	if err != nil {
		return err
	}

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()
	defer output.Body.Close()

	buf := make([]byte, BufferSize)
	eof := false
	for !eof {
		read, readErr := output.Body.Read(buf)
		if readErr != nil && readErr != io.EOF {
			return readErr
		}
		written, writeErr := file.Write(buf[:read])
		if writeErr != nil {
			return writeErr
		}
		if written != read {
			return errors.New("mismatch between read and write len")
		}

		eof = readErr == io.EOF
	}

	return nil
}
