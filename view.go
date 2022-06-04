package main

import (
	"github.com/nsf/termbox-go"
)

type ObjectListView struct {
	cursor        Cursor
	width, height int

	client     Storage
	bucket     string
	prefix     string
	objectList *ObjectList
}

func NewObjectListView() *ObjectListView {
	return &ObjectListView{
		cursor: Cursor{},
		client: NewStorage(AWS),
	}
}

func (v *ObjectListView) Start() (err error) {
	v.width, v.height = termbox.Size()

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

func (v *ObjectListView) printObjectList() {
	for _, obj := range v.objectList.List {
		v.PrintObject(obj)
	}
	v.cursor.Reset()
	v.setActiveLine()
	termbox.Flush()
}

func (v *ObjectListView) PrintObject(obj Object) {
	switch obj.Type {
	case BucketType:
		v.printBucketLine(obj)
	case DirectoryType:
		v.printDirectoryLine(obj)
	case FileType:
		v.printFileLine(obj)
	}
	v.cursor.Down()
}

func (v *ObjectListView) printBucketLine(obj Object) {
	TBPrintMsg(v.cursor.x, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, "\U0001FAA3")
	TBPrintMsg(v.cursor.x+3, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, obj.Key)
}

func (v *ObjectListView) printDirectoryLine(obj Object) {
	TBPrintMsg(v.cursor.x, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, "\U0001FC41")
	TBPrintMsg(v.cursor.x+3, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, obj.Key)
}

func (v *ObjectListView) printFileLine(obj Object) {
	TBPrintMsg(v.cursor.x, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, obj.Key)
}

func (v *ObjectListView) setActiveLine() {
	for i := 0; i < v.width; i++ {
		termbox.SetBg(v.cursor.x+i, v.cursor.y, termbox.ColorWhite)
		termbox.SetFg(v.cursor.x+i, v.cursor.y, termbox.ColorBlack)
	}
}

func (v *ObjectListView) setDefaultLine() {
	for i := 0; i < v.width; i++ {
		termbox.SetBg(v.cursor.x+i, v.cursor.y, termbox.ColorDefault)
		termbox.SetFg(v.cursor.x+i, v.cursor.y, termbox.ColorWhite)
	}
}

func (v *ObjectListView) Up() {
	if v.cursor.isTop() {
		return
	}
	v.setDefaultLine()
	v.cursor.Up()
	v.setActiveLine()
	termbox.Flush()
}

func (v *ObjectListView) Down() {
	if v.cursor.isBottom() {
		return
	}
	v.setDefaultLine()
	v.cursor.Down()
	v.setActiveLine()
	termbox.Flush()
}
