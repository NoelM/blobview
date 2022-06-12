package main

type Cloud int

const (
	AWS Cloud = iota
)

var Protocol = map[Cloud]string{
	AWS: "s3",
}

const BufferSize = 1024 * 1024

type Storage interface {
	Start() error

	ListBuckets() (*ObjectList, error)

	ListObjects(bucket, prefix string) (*ObjectList, error)
	ListObjectsNext(bucket, prefix, token string) (*ObjectList, error)

	Download(object Object, destination string) error
}

func NewStorage(cloud Cloud) Storage {
	switch cloud {
	case AWS:
		return NewAWSS3Driver()
	default:
		return nil
	}
}
