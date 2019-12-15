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
