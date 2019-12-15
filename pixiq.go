package pixiq

// NewImages returns a factory of images which can be passed around to every place where you construct new images.
func NewImages() *Images {
	return &Images{}
}

// Images is a factory of images used to create new images.
type Images struct {
}

// New creates an Image with specified size given in pixels.
func (i *Images) New(width, height int) *Image {
	var w, h int
	if width > 0 {
		w = width
	}
	if height > 0 {
		h = height
	}
	return &Image{
		width:  w,
		height: h,
	}
}

// Image is a 2D picture composed of pixels each having specific color. The cost of creating an Image is huge therefore
// new images should be created sporadically, ideally when the application starts.
type Image struct {
	width, height int
}

// Width returns the number of pixels in a row.
func (i *Image) Width() int {
	return i.width
}

// Height returns the number of pixels in a column.
func (i *Image) Height() int {
	return i.height
}

// Selection makes a rectangular selection starting at a given position. The position has to be top-left position.
// Both x and y can be negative, meaning that selection starts outside the image.
func (i *Image) Selection(x int, y int) Selection {
	return Selection{
		x: x,
		y: y,
	}
}

// WholeImageSelection make selection of entire image
func (i *Image) WholeImageSelection() Selection {
	return Selection{
		width:  i.width,
		height: i.height,
	}
}

// Selection marks a selection on top of the image.
type Selection struct {
	x, y, width, height int
}

// Width returns the width of selection in pixels.
func (s Selection) Width() int {
	return s.width
}

// Height returns the height of selection in pixels.
func (s Selection) Height() int {
	return s.height
}

// X returns the starting position
func (s Selection) X() int {
	return s.x
}

// Y returns the starting position
func (s Selection) Y() int {
	return s.y
}

// WithSize creates a new selection with specified size in pixels. Width or height are clamped to 0
func (s Selection) WithSize(width, height int) Selection {
	if width > 0 {
		s.width = width
	} else {
		s.width = 0
	}
	if height > 0 {
		s.height = height
	} else {
		s.height = 0
	}
	return s
}
