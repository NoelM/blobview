package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func TBPrintMsg(x, y int, fg, bg termbox.Attribute, msg string) int {
	size := 0
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
		size++
	}
	return size
}

func IntMin(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func IntMax(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func printPath(obj Object) string {
	return fmt.Sprintf("%s://%s", Protocol[obj.Provider], obj.Key)
}

func clearLine(c *Cursor, fg, bg termbox.Attribute) {
	c.LineOrigin()
	for !c.IsRight() {
		termbox.SetCell(c.x, c.y, ' ', fg, bg)
		c.Right()
	}
}
