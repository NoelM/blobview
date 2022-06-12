package main

type Cursor struct {
	x, y             int
	xOrigin, yOrigin int
	xSize, ySize     int
	width, height    int
}

func (c *Cursor) IsTop() bool {
	return c.y == c.yOrigin
}

func (c *Cursor) Up() {
	if c.IsTop() {
		return
	}
	c.y--
}

func (c *Cursor) MoveUp(i int) {
	c.y = IntMax(c.yOrigin, c.y-i)
}

func (c *Cursor) IsBottom() bool {
	return c.y == c.yOrigin+c.ySize-1
}

func (c *Cursor) Down() {
	c.y++
}

func (c *Cursor) MoveDown(i int) {
	c.y = IntMin(c.yOrigin+c.ySize-1, c.y+i)
}

func (c *Cursor) IsLeft() bool {
	return c.x == c.xOrigin
}

func (c *Cursor) Left() {
	c.x--
}

func (c *Cursor) MoveLeft(i int) {
	c.x = IntMax(c.xOrigin, c.x-i)
}

func (c *Cursor) IsRight() bool {
	return c.x == c.xOrigin+c.xSize-1
}

func (c *Cursor) Right() {
	c.x++
}

func (c *Cursor) MoveRight(i int) {
	c.x = IntMin(c.xOrigin+c.xSize-1, c.x+i)
}

func (c *Cursor) LineOrigin() {
	c.x = c.xOrigin
}

func (c *Cursor) NextLine() {
	c.x = c.xOrigin
	if c.IsBottom() {
		return
	}
	c.y++
}

func (c *Cursor) PrevLine() {
	c.x = c.xOrigin
	if c.IsTop() {
		return
	}
	c.y--
}

func (c *Cursor) Sync(width, height int) {
	c.width, c.height = width, height
}

func (c *Cursor) Reset() {
	x, y := c.xOrigin, c.yOrigin

	if c.xOrigin < 0 {
		x = c.width - 1 + c.xOrigin
	}
	if c.yOrigin < 0 {
		y = c.height - 1 + c.yOrigin
	}

	c.x, c.y = x, y
}
