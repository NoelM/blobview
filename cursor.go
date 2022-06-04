package main

import "github.com/nsf/termbox-go"

type Cursor struct {
	x, y          int
	height, width int
}

func (c *Cursor) isTop() bool {
	return c.y == 0
}

func (c *Cursor) isBottom() bool {
	c.width, c.height = termbox.Size()
	return c.y == c.height
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

func (c *Cursor) isLeftSide() bool {
	return c.x == 0
}

func (c *Cursor) Left() {
	if c.isLeftSide() {
		return
	}
	c.x--
}
func (c *Cursor) isRightSide() bool {
	c.width, c.height = termbox.Size()
	return c.x == c.width
}

func (c *Cursor) Right() {
	if c.isRightSide() {
		return
	}
	c.x++
}

func (c *Cursor) Reset() {
	c.x, c.y = 0, 0
}
