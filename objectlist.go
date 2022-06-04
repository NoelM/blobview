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
	Active int

	List []Object
}

func (o *ObjectList) ActiveUp() {
	if o.Active == 0 {
		return
	}
	o.Active--
}

func (o *ObjectList) ActiveDown() {
	if o.Active == len(o.List)-1 {
		return
	}
	o.Active++
}

func (o *ObjectList) GetActiveObject() Object {
	return o.List[o.Active]
}
