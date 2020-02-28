package image

// AcceleratedImage is an image processed by external device (outside the CPU).
// The mentioned device might be a video card.
type AcceleratedImage interface {
	// Upload transfers pixels from RAM to external memory (such as VRAM).
	//
	// Pixels must have pixel colors sorted by coordinates.
	// Pixels are send for last line first, from left to right.
	// Pixels slice holds all image pixels and therefore must have size width*height
	//
	// Implementations must not retain pixels slice and make a copy instead.
	Upload(pixels []Color)
	// Download transfers pixels from external memory (such as VRAM) to RAM.
	//
	// Output will have pixel colors sorted by coordinates.
	// Pixels are send for last line first, from left to right.
	// Output must be of size width*height.
	//
	// If the image has not been uploaded before then Download should fill
	// output with Transparent color.
	//
	// Implementations must not retain output.
	Download(output []Color)
}

// New creates an Image with specified size given in pixels.
// Will panic if AcceleratedImage is nil or width and height are negative
func New(width, height int, acceleratedImage AcceleratedImage) *Image {
	if acceleratedImage == nil {
		panic("nil acceleratedImage")
	}
	if width < 0 {
		panic("negative width")
	}
	if height < 0 {
		panic("negative height")
	}
	return &Image{
		width:            width,
		height:           height,
		heightMinusOne:   height - 1,
		pixels:           make([]Color, width*height),
		acceleratedImage: acceleratedImage,
		selectionsCache:  make([]AcceleratedImageSelection, 0, 4),
	}
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
	heightMinusOne   int
	pixels           []Color
	acceleratedImage AcceleratedImage
	selectionsCache  []AcceleratedImageSelection
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
	y := s.image.heightMinusOne - localY - s.y
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
	y := s.image.heightMinusOne - localY - s.y
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

// AcceleratedImageLocation is a location of a AcceleratedImage
type AcceleratedImageLocation struct {
	X, Y, Width, Height int
}

// AcceleratedImageSelection is same for AcceleratedImage as Selection for *Image
type AcceleratedImageSelection struct {
	Location AcceleratedImageLocation
	Image    AcceleratedImage
}

// AcceleratedCommand is a command executed externally (outside the CPU).
type AcceleratedCommand interface {
	// Run should put the results into the output selection of AcceleratedImage,
	// so that next time AcceleratedImage.Download is called modified pixels are
	// downloaded.
	//
	// Run might return error when output or selections cannot be used. Usually
	// the reason for that is they were not created in a given context (such
	// as OpenGL context).
	//
	// Implementations must not retain selections.
	Run(output AcceleratedImageSelection, selections []AcceleratedImageSelection)
}

// Modify runs the AcceleratedCommand and put results into the Selection.
// This method ensures that all passed selections are uploaded before the command
// is called. Selections get converted into AcceleratedImageSelection and
// passed to the command.Run.
func (s Selection) Modify(command AcceleratedCommand, selections ...Selection) {
	if command == nil {
		return
	}
	convertedSelections := s.image.selectionsCache[:0]
	for _, selection := range selections {
		selection.image.acceleratedImage.Upload(selection.image.pixels)
		convertedSelections = append(convertedSelections, selection.toAcceleratedImageSelection())
	}
	command.Run(s.toAcceleratedImageSelection(), convertedSelections)
	s.image.acceleratedImage.Download(s.image.pixels)
}

func (s Selection) toAcceleratedImageSelection() AcceleratedImageSelection {
	return AcceleratedImageSelection{
		Location: AcceleratedImageLocation{
			X:      s.x,
			Y:      s.y,
			Width:  s.width,
			Height: s.height,
		},
		Image: s.image.acceleratedImage,
	}
}
