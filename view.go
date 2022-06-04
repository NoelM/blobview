package main

import (
	"github.com/nsf/termbox-go"
)

const ScreenSize = 50
const ScreenWidth = 80

type Cursor struct {
	x, y int
}

func (c *Cursor) isTop() bool {
	return c.y == 0
}

func (c *Cursor) isBottom() bool {
	return c.y == ScreenSize
}

func (c *Cursor) Up() {
	if c.isTop() {
		return
	}
	c.y--
}

func (c *Cursor) Down() {
	if c.isBottom() {
		return
	}
	c.y++
}

func (c *Cursor) Reset() {
	c.x, c.y = 0, 0
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

func (v *View) PrintObjectList(objects []Object) {
	for _, obj := range objects {
		v.PrintObject(obj)
	}
	v.cursor.Reset()
	v.setActiveLine()
	termbox.Flush()
}

func (v *View) PrintObject(obj Object) {
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

func (v *View) printBucketLine(obj Object) {
	TBPrintMsg(v.cursor.x, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, "\U0001FAA3")
	TBPrintMsg(v.cursor.x+3, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, obj.Name)
}

func (v *View) printDirectoryLine(obj Object) {
	TBPrintMsg(v.cursor.x, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, "\U0001FC41")
	TBPrintMsg(v.cursor.x+3, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, obj.Name)
}

func (v *View) printFileLine(obj Object) {
	TBPrintMsg(v.cursor.x, v.cursor.y, termbox.ColorWhite, termbox.ColorDefault, obj.Name)
}

func (v *View) setActiveLine() {
	for i := 0; i < ScreenWidth; i++ {
		termbox.SetBg(v.cursor.x+i, v.cursor.y, termbox.ColorWhite)
		termbox.SetFg(v.cursor.x+i, v.cursor.y, termbox.ColorBlack)
	}
}

func (v *View) setDefaultLine() {
	for i := 0; i < ScreenWidth; i++ {
		termbox.SetBg(v.cursor.x+i, v.cursor.y, termbox.ColorDefault)
		termbox.SetFg(v.cursor.x+i, v.cursor.y, termbox.ColorWhite)
	}
}

func (v *View) Up() {
	if v.cursor.isTop() {
		return
	}
	v.setDefaultLine()
	v.cursor.Up()
	v.setActiveLine()
	termbox.Flush()
}

func (v *View) Down() {
	if v.cursor.isBottom() {
		return
	}
	v.setDefaultLine()
	v.cursor.Down()
	v.setActiveLine()
	termbox.Flush()
}
