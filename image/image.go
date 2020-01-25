package image

import "errors"

// AcceleratedImage is an image processed externally (outside the CPU).
type AcceleratedImage interface {
	// Upload send pixels colors sorted by coordinates.
	// First all pixels are sent for y=0, from left to right.
	Upload(pixels []Color)
	// Downloads pixels by filling output Color slice
	Download(output []Color)
	// Create a modification program using AcceleratedProgram. The results should
	// be store in in a given selection
	// Passed AcceleratedImageSelection is always clamped to image boundaries
	Modify(AcceleratedProgram, AcceleratedImageSelection) (AcceleratedModification, error)
}

type AcceleratedModification interface {
}

type AcceleratedImageSelection struct {
	X, Y, Width, Height int
}

// AcceleratedProgram is a program executed externally (outside the CPU).
type AcceleratedProgram interface {
}

// New creates an Image with specified size given in pixels.
// Will panic if AcceleratedImage is nil
// Will return error if width and height are negative
func New(width, height int, acceleratedImage AcceleratedImage) (*Image, error) {
	if acceleratedImage == nil {
		panic("nil acceleratedImage")
	}
	if width < 0 {
		return nil, errors.New("negative width")
	}
	if height < 0 {
		return nil, errors.New("negative height")
	}
	return &Image{
		width:            width,
		height:           height,
		pixels:           make([]Color, width*height),
		acceleratedImage: acceleratedImage,
	}, nil
}

// Image is a 2D picture composed of pixels each having a specific color.
// Image is using 2 coordinates: X and Y to specify the position of a pixel.
// The origin (0,0) is at the top-left corner of the image.
//
// The cost of creating an Image is huge therefore new images should be created
// sporadically, ideally when the application starts.
type Image struct {
	width            int
	height           int
	pixels           []Color
	acceleratedImage AcceleratedImage
}

// Width returns the number of pixels in a row.
func (i *Image) Width() int {
	return i.width
}

// Height returns the number of pixels in a column.
func (i *Image) Height() int {
	return i.height
}

// Selection creates an area pointing to the image at a given starting position
// (x and y). The position must be a top left corner of the selection.
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

// Upload uploads all image pixels to associated AcceleratedImage.
// This method should be called rarely. Image pixels are uploaded automatically
// when needed.
//
// DEPRECATED - this method will be removed in next release
func (i *Image) Upload() {
	i.acceleratedImage.Upload(i.pixels)
}

// Selection points to a specific area of the image. It has a starting position
// (top-left corner) and optional size. Most Selection methods - such as Color,
// SetColor and Selection use local coordinates as parameters. Top-left corner
// of selection has (0,0) local coordinates.
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

// WithSize creates a new selection with specified size in pixels.
// Negative width or height are constrained to 0.
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

// Selection makes a new selection using the coordinates of existing selection.
// Passed coordinates are local, which means that the top-left corner of existing
// selection is equivalent to localX=0, localY=0. Both coordinates can be negative,
// meaning that selection starts outside the original selection.
func (s Selection) Selection(localX, localY int) Selection {
	return Selection{
		x:     localX + s.x,
		y:     localY + s.y,
		image: s.image,
	}
}

// Color returns the color of the pixel at a specific position.
// Passed coordinates are local, which means that the top-left corner of selection
// is equivalent to localX=0, localY=0. Negative coordinates are supported.
// If pixel is outside the image boundaries then transparent color is returned.
// It is possible to get the color outside the selection.
func (s Selection) Color(localX, localY int) Color {
	x := localX + s.x
	if x < 0 {
		return Transparent
	}
	y := localY + s.y
	if y < 0 {
		return Transparent
	}
	if x >= s.image.width {
		return Transparent
	}
	index := x + y*s.image.width
	if len(s.image.pixels) <= index {
		return Transparent
	}
	return s.image.pixels[index]
}

// SetColor sets the color of the pixel at specific position.
// Passed coordinates are local, which means that the top-left corner of selection
// is equivalent to localX=0, localY=0. Negative coordinates are supported.
// If pixel is outside the image boundaries then nothing happens.
// It is possible to set the color outside the selection.
func (s Selection) SetColor(localX, localY int, color Color) {
	x := localX + s.x
	if x < 0 {
		return
	}
	y := localY + s.y
	if y < 0 {
		return
	}
	if x >= s.image.width {
		return
	}
	index := x + y*s.image.width
	if len(s.image.pixels) <= index {
		return
	}
	s.image.pixels[index] = color
}

func (s Selection) toAcceleratedImageSelection() AcceleratedImageSelection {
	var (
		x      = s.x
		y      = s.y
		width  = s.width
		height = s.height
	)
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= s.image.width {
		x = 0
	}
	if x+width > s.image.width {
		width = s.image.width - x
	}
	if y >= s.image.height {
		y = 0
	}
	if y+height > s.image.height {
		height = s.image.height - y
	}
	return AcceleratedImageSelection{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

type SelectionModification struct {
}

func (s Selection) Modify(acceleratedProgram AcceleratedProgram, cpuProgram func(SelectionModification)) error {
	if acceleratedProgram == nil {
		return errors.New("nil acceleratedProgram")
	}
	if cpuProgram == nil {
		return errors.New("nil cpuProgram")
	}

	_, err := s.image.acceleratedImage.Modify(acceleratedProgram, s.toAcceleratedImageSelection())
	if err != nil {
		return err
	}
	cpuProgram(SelectionModification{})
	return nil
}
