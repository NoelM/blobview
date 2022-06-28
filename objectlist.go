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
	Provider Cloud
	Bucket   string
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

func (o *ObjectList) SetActive(id int) {
	o.Active = id
}

func (o *ObjectList) IsFirstActive() bool {
	return o.Active == 0
}

func (o *ObjectList) ActiveUp() {
	if o.IsFirstActive() {
		return
	}
	o.Active--
}

func (o *ObjectList) IsLastActive() bool {
	return o.Active == len(o.List)-1
}

func (o *ObjectList) ActiveDown() {
	if o.IsLastActive() {
		return
	}
	o.Active++
}

func (o *ObjectList) GetActiveObject() Object {
	return o.List[o.Active]
}

func (o *ObjectList) HasPrevious() bool {
	return !o.IsFirstActive()
}

func (o *ObjectList) HasNext() bool {
	return !o.IsLastActive() || len(o.Token) != 0
}
