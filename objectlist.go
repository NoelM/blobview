package main

import (
	"time"
)

type ObjectType int

const (
	BucketType ObjectType = iota
	DirectoryType
	FileType
)

type Object struct {
	Key   string
	Type  ObjectType
	Bytes int64
	Date  time.Time
}

type ObjectList struct {
	Bucket string
	Prefix string
	Token  string

	List []Object
}
