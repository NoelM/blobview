package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"path"
	"strings"
	"time"
)

type ObjectListView struct {
	listCursor    Cursor
	headerCursor  Cursor
	width, height int
	maxKeySize    int
	columnFormat  string

	client     Storage
	objectList *ObjectList
}

func NewObjectListView() *ObjectListView {
	return &ObjectListView{
		headerCursor: Cursor{
			xOrigin: 1,
			yOrigin: 0,
		},
		listCursor: Cursor{
			xOrigin: 1,
			yOrigin: 3,
		},
		client: NewStorage(AWS),
	}
}

func (v *ObjectListView) Start() (err error) {
	v.width, v.height = termbox.Size()
	v.headerCursor.xSize = v.width - 2
	v.headerCursor.ySize = 3
	v.headerCursor.Reset()

	v.listCursor.xSize = v.width - 2
	v.listCursor.ySize = v.height - 3
	v.listCursor.Reset()

	if err = v.client.Start(); err != nil {
		return err
	}

	v.objectList, err = v.client.ListBuckets()
	if err != nil {
		return err
	}

	v.printObjectList()
	return nil
}

func (v *ObjectListView) Dive() {
	active := v.objectList.GetActiveObject()
	if active.Type == FileType {
		return
	}

	var objectList *ObjectList
	var err error

	if active.Type == DirectoryType {
		objectList, err = v.client.ListObjects(v.objectList.Bucket, active.Key)
	} else if active.Type == BucketType {
		objectList, err = v.client.ListObjects(active.Key, "")
	}

	if err != nil {
		return
	}

	v.objectList = objectList

	v.Reset()
	v.printObjectList()
}

func (v *ObjectListView) Back() {
	if v.objectList.Bucket == "" {
		return
	}

	if v.objectList.Prefix == "" {
		objectList, err := v.client.ListBuckets()
		if err != nil {
			return
		}

		v.objectList = objectList

		v.Reset()
		v.printObjectList()
	} else {
		prefixList := strings.Split(v.objectList.Prefix, "/")
		var newPrefix string
		if len(prefixList) > 2 {
			newPrefix = path.Join(prefixList[:len(prefixList)-2]...) + "/"
		}

		objectList, err := v.client.ListObjects(v.objectList.Bucket, newPrefix)
		if err != nil {
			return
		}

		v.objectList = objectList

		v.Reset()
		v.printObjectList()
	}
}

func (v *ObjectListView) Reset() {
	v.listCursor.Reset()
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
}

func (v *ObjectListView) printObjectList() {
	v.headerCursor.Reset()
	v.printHeaders()

	for _, obj := range v.objectList.List {
		v.printObject(obj)
	}

	v.listCursor.Reset()
	v.setActiveLine()

	termbox.Flush()
}

func (v *ObjectListView) printObject(obj Object) {
	switch obj.Type {
	case BucketType:
		v.printBucketLine(obj)
	case DirectoryType:
		v.printDirectoryLine(obj)
	case FileType:
		v.printFileLine(obj)
	}
	v.listCursor.NextLine()
}

func (v *ObjectListView) printBucketLine(obj Object) {
	line := fmt.Sprintf(v.columnFormat, "\U0001FAA3", obj.PrintKey, obj.Date.Format(time.RFC822Z), "")
	v.listCursor.LineOrigin()

	TBPrintMsg(v.listCursor.x, v.listCursor.y, termbox.ColorWhite, termbox.ColorDefault, line)
}

func (v *ObjectListView) printDirectoryLine(obj Object) {
	line := fmt.Sprintf(v.columnFormat, "\U0001F4C1", obj.PrintKey, "", "")
	v.listCursor.LineOrigin()

	TBPrintMsg(v.listCursor.x, v.listCursor.y, termbox.ColorWhite, termbox.ColorDefault, line)
}

func (v *ObjectListView) printFileLine(obj Object) {
	if obj.PrintKey == "" {
		return
	}
	line := fmt.Sprintf(v.columnFormat, "\U0001F4C4", obj.PrintKey, obj.Date.Format(time.RFC822Z), fmt.Sprintf("%5d", obj.Bytes/1024/1024))
	v.listCursor.LineOrigin()

	TBPrintMsg(v.listCursor.x, v.listCursor.y, termbox.ColorWhite, termbox.ColorDefault, line)
}

func (v *ObjectListView) setActiveLine() {
	v.listCursor.LineOrigin()

	for !v.listCursor.IsRight() {
		termbox.SetBg(v.listCursor.x, v.listCursor.y, termbox.ColorWhite)
		termbox.SetFg(v.listCursor.x, v.listCursor.y, termbox.ColorBlack)

		v.listCursor.Right()
	}
}

func (v *ObjectListView) setDefaultLine() {
	v.listCursor.LineOrigin()

	for !v.listCursor.IsRight() {
		termbox.SetBg(v.listCursor.x, v.listCursor.y, termbox.ColorDefault)
		termbox.SetFg(v.listCursor.x, v.listCursor.y, termbox.ColorWhite)

		v.listCursor.Right()
	}
}

func (v *ObjectListView) Up() {
	v.objectList.ActiveUp()

	if v.listCursor.IsTop() {
		return
	}
	v.setDefaultLine()
	v.listCursor.Up()
	v.setActiveLine()

	termbox.Flush()
}

func (v *ObjectListView) Down() {
	v.objectList.ActiveDown()

	if v.listCursor.IsBottom() {
		return
	}
	v.setDefaultLine()
	v.listCursor.Down()
	v.setActiveLine()

	termbox.Flush()
}

func (v *ObjectListView) printHeaders() {
	if v.objectList.Bucket == "" {
		v.printBucketListHeaders()
	} else {
		v.printObjectListHeaders()
	}
	v.headerCursor.NextLine()
	v.printColumnHeaders()
}

func (v *ObjectListView) printBucketListHeaders() {
	n := TBPrintMsg(v.headerCursor.x, v.headerCursor.y, termbox.ColorWhite, termbox.ColorBlue, "== BLOBVIEW 0.1 == Bucket List")
	v.headerCursor.MoveRight(n)

	for !v.headerCursor.IsRight() {
		termbox.SetBg(v.headerCursor.x, v.headerCursor.y, termbox.ColorBlue)
		v.headerCursor.Right()
	}
}

func (v *ObjectListView) printObjectListHeaders() {
	n := TBPrintMsg(v.headerCursor.x, v.headerCursor.y, termbox.ColorWhite, termbox.ColorBlue, "== BLOBVIEW 0.1 == s3://"+v.objectList.Bucket+"/"+v.objectList.Prefix)
	v.headerCursor.MoveRight(n)

	for !v.headerCursor.IsRight() {
		termbox.SetBg(v.headerCursor.x, v.headerCursor.y, termbox.ColorBlue)
		v.headerCursor.Right()
	}
}

func (v *ObjectListView) printColumnHeaders() {
	maxPrintPrefixSize := 0
	for _, obj := range v.objectList.List {
		maxPrintPrefixSize = IntMax(maxPrintPrefixSize, len(obj.PrintKey))
	}

	maxNameSize := IntMax(maxPrintPrefixSize, 20)
	maxNameSize = IntMin(maxNameSize, v.width-25)

	v.columnFormat = fmt.Sprintf("%%-2s %%-%ds %%-24s %%-5s", maxNameSize)

	line := fmt.Sprintf(v.columnFormat, "", "NAME", "DATE", "SIZE")
	v.headerCursor.LineOrigin()

	n := TBPrintMsg(v.headerCursor.x, v.headerCursor.y, termbox.ColorWhite, termbox.ColorBlue, line)
	v.headerCursor.MoveRight(n)

	for !v.headerCursor.IsRight() {
		termbox.SetBg(v.headerCursor.x, v.headerCursor.y, termbox.ColorBlue)
		v.headerCursor.Right()
	}
}
