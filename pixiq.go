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
		pixels: make([]Color, w*h),
	}
}

// Image is a 2D picture composed of pixels each having a specific color. Image is using 2 coordinates: X and Y to
// specify the position of a pixel. The origin (0,0) is at the top-left corner of the image.
//
// The cost of creating an Image is huge therefore new images should be created sporadically, ideally when
// the application starts.
type Image struct {
	width  int
	height int
	pixels []Color
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
		x:     x,
		y:     y,
		image: i,
	}
}

// WholeImageSelection make selection of entire image.
func (i *Image) WholeImageSelection() Selection {
	return i.Selection(0, 0).WithSize(i.width, i.height)
}

// Selection marks a selection on top of the image. Most Selection methods - such as Color, SetColor and Selection use local
// coordinates which means that top-left corner of Selection is an origin (0,0). This position can be different than
// position given in image coordinates.
type Selection struct {
	image         *Image
	x, y          int
	width, height int
}

// Image returns image for which the selection was made.
func (s Selection) Image() *Image {
	return s.image
}

// Width returns the width of selection in pixels.
func (s Selection) Width() int {
	return s.width
}

// Height returns the height of selection in pixels.
func (s Selection) Height() int {
	return s.height
}

// ImageX returns the starting position in image coordinates.
func (s Selection) ImageX() int {
	return s.x
}

// ImageY returns the starting position in image coordinates.
func (s Selection) ImageY() int {
	return s.y
}

// WithSize creates a new selection with specified size in pixels. Width or height are clamped to 0 if necessarily.
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

// Selection makes a new selection using the coordinates of existing selection. Passed coordinates are local,
// which means that the top-left corner of existing selection is equivalent to localX=0, localY=0.
// Both coordinates can be negative, meaning that selection starts outside the original selection.
func (s Selection) Selection(localX, localY int) Selection {
	return Selection{
		x:     localX + s.x,
		y:     localY + s.y,
		image: s.image,
	}
}

// Color returns the color of the pixel at a specific position. Passed coordinates are local,
// which means that the top-left corner of selection is equivalent to localX=0, localY=0.
// Negative coordinates are supported. If pixel is outside the image boundaries then transparent color is returned.
func (s Selection) Color(localX, localY int) Color {
	if localX < 0 || localY < 0 || localY >= s.image.height {
		return Transparent
	}
	return s.image.pixels[localY*s.image.width+localX]
}

// SetColor set the color of the pixel at specific position. Passed coordinates are local,
// which means that the top-left corner of selection is equivalent to localX=0, localY=0.
// Negative coordinates are supported. If pixel is outside the image boundaries then nothing happens.
func (s Selection) SetColor(localX, localY int, color Color) {
	s.image.pixels[0] = color
}
