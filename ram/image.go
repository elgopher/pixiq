package ram

import "github.com/jacekolszak/pixiq/color"

func NewImage(width, height int) *Image {
	return &Image{
		width:  width,
		height: height,
	}
}

type Image struct {
	width, height int
}

// Width returns width of the image given at the creation time
func (i *Image) Width() int {
	return i.width
}

// Height returns width of the image given at the creation time
func (i *Image) Height() int {
	return i.height
}

// Selection creates a rectangle selection with given boundaries.
func (i *Image) Selection(x, y, width, height int) Selection {
	return Selection{
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

type Selection struct {
	x, y, width, height int
}

func (s Selection) X() int {
	return s.x
}

func (s Selection) Y() int {
	return s.y
}

func (s Selection) Width() int {
	return s.width
}

func (s Selection) Height() int {
	return s.height
}

// Line returns line with given local number
func (s Selection) Line(y int) Line {
	return Line{}
}

type Line struct {
}

// Set sets pixel color with given local X coordinate
func (l Line) Set(x int, c color.Color) {
}

func (l Line) Get(x int) color.Color {
	return color.Color{}
}
