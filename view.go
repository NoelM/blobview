package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

const ScreenSize = 50

type Cursor struct {
	x int
	y int
}

func (c *Cursor) Up() {
	if c.y == 0 {
		return
	}
	c.y--
}

func (c *Cursor) Down() {
	if c.y == ScreenSize {
		return
	}
	c.y++
}

type View struct {
	cursor Cursor
	bucket string
	path   string
}

func NewView() *View {

	return &View{
		cursor: Cursor{0, 0},
	}
}

type Object struct {
	Name      string
	Bucket    bool
	Directory bool
	Bytes     int64
	Date      time.Time
}

func (v *View) PrintObjectList(objects []Object) {
	for _, obj := range objects {
		v.PrintObject(obj)
	}
	termbox.Flush()
}

func (v *View) PrintObject(obj Object) {
	TBPrintMsg(v.cursor.x, v.cursor.y, termbox.ColorWhite, termbox.ColorBlack, obj.Name)
	v.cursor.Down()
}
