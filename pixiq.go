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

func (i *Image) Selection(x int, y int) Selection {
	var s Selection
	s.x = x
	s.y = y
	return s
}

type Selection struct {
	x, y, width, height int
}

func (s Selection) Width() int {
	return s.width
}

func (s Selection) Height() int {
	return s.height
}

func (s Selection) X() int {
	return s.x
}

func (s Selection) Y() int {
	return s.y
}

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
