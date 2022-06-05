package main

import (
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
