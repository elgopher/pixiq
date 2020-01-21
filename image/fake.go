package image

// NewFakeImages creates an instance of FakeImages
func NewFakeImages() *FakeImages {
	return &FakeImages{
		registeredCalls: map[interface{}]FakeCall{},
	}
}

// FakeImages is a factory of Fake images and registry of calls which can
// be executed on them. It is useful in unit testing.
type FakeImages struct {
	registeredCalls map[interface{}]FakeCall
}

// NewAcceleratedImage creates a new FakeAcceleratedImage instance which can
// be used in unit testing.
// TODO Width and height are constrained to zero if negative.
func (i *FakeImages) NewAcceleratedImage(width, height int) *FakeAcceleratedImage {
	return &FakeAcceleratedImage{
		registeredCalls: i.registeredCalls,
		width:           width,
		height:          height,
		pixels:          make([]Color, width*height),
	}
}

// NewImageWithFakeAcceleration can be used in unit tests for testing
// CPU-based functionality
func (i *FakeImages) NewImageWithFakeAcceleration(width, height int) *Image {
	return New(width, height, i.NewAcceleratedImage(width, height))
}

// RegisterCall registers a FakeCall implementation which then can be used in Modify.
func (i *FakeImages) RegisterCall(call FakeCall) {
	if call == nil {
		panic("nil call")
	}
	i.registeredCalls[call] = call
}

// FakeCall is an interface which must be implemented by any fake calls
// used in FakeAcceleratedImage.Modify().
type FakeCall interface {
	Run(selection AcceleratedFragmentLocation, image *FakeAcceleratedImage)
}

// AddColor creates and registers an
func (i *FakeImages) AddColor(c Color) AcceleratedCall {
	call := &fakeAddColor{color: c}
	i.registeredCalls[call] = call
	return call
}

type fakeAddColor struct {
	color Color
}

func (f *fakeAddColor) Run(selection AcceleratedFragmentLocation, image *FakeAcceleratedImage) {
	for x := selection.X; x < selection.X+selection.Width; x++ {
		for y := selection.Y; y < selection.Y+selection.Height; y++ {
			index := x + y*image.width
			image.pixels[index] = RGBA(
				image.pixels[index].R()+f.color.R(),
				image.pixels[index].G()+f.color.G(),
				image.pixels[index].B()+f.color.B(),
				image.pixels[index].A()+f.color.A(),
			)
		}
	}
}

// FakeAcceleratedImage is an AcceleratedImage which emulates the behaviour of
// a real AcceleratedImage and can be used in unit tests.
// FakeAcceleratedImage stores the uploaded pixels in RAM.
type FakeAcceleratedImage struct {
	registeredCalls map[interface{}]FakeCall
	pixels          []Color
	width           int
	height          int
}

// Upload copies the pixels to RAM.
func (i *FakeAcceleratedImage) Upload(input AcceleratedFragmentPixels) {
	inputOffset := input.StartingPosition
	location := input.Location
	for y := 0; y < location.Height; y++ {
		for x := 0; x < location.Width; x++ {
			index := y*i.width + x + location.X + location.Y*i.width
			i.pixels[index] = input.Pixels[inputOffset]
			inputOffset++
		}
		inputOffset += input.Stride - location.Width
	}
}

// Download copies the pixels from RAM to AcceleratedFragmentPixels struct.
func (i *FakeAcceleratedImage) Download(output AcceleratedFragmentPixels) {
	location := output.Location
	outputOffset := output.StartingPosition
	for y := 0; y < location.Height; y++ {
		for x := 0; x < location.Width; x++ {
			index := y*i.width + x + location.X + location.Y*i.width
			output.Pixels[outputOffset] = i.pixels[index]
			outputOffset++
		}
		outputOffset += output.Stride - location.Width
	}
}

// Modify executes the call. Call must be a FakeCall implementation registered
// before using the FakeAcceleratedImage.RegisterCall method.
func (i *FakeAcceleratedImage) Modify(selection AcceleratedFragmentLocation, call AcceleratedCall) {
	if call == nil {
		panic("nil call")
	}
	fakeCall, ok := i.registeredCalls[call]
	if !ok {
		panic("invalid call")
	}
	fakeCall.Run(selection, i)
}
