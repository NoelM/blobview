package main

type Cloud int

const (
	AWS Cloud = iota
)

type Storage interface {
	Start() error

	ListBuckets() (*ObjectList, error)

	ListObjects(bucket, prefix string) (*ObjectList, error)
	ListObjectsNext(bucket, prefix, token string) (*ObjectList, error)
}

func NewStorage(cloud Cloud) Storage {
	switch cloud {
	case AWS:
		return NewAWSS3Driver()
	default:
		return nil
	}
}
