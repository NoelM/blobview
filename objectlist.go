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
	Key      string
	PrintKey string
	Type     ObjectType
	Bytes    int64
	Date     time.Time
}

type ObjectList struct {
	Bucket string
	Prefix string
	Token  string
	Active int

	List []Object
}

func (o *ObjectList) FirstActive() bool {
	return o.Active == 0
}

func (o *ObjectList) ActiveUp() {
	if o.FirstActive() {
		return
	}
	o.Active--
}

func (o *ObjectList) LastActive() bool {
	return o.Active == len(o.List)-1
}

func (o *ObjectList) ActiveDown() {
	if o.LastActive() {
		return
	}
	o.Active++
}

func (o *ObjectList) GetActiveObject() Object {
	return o.List[o.Active]
}
