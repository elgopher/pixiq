package image

func NewFake() *Fake {
	return &Fake{calls: map[interface{}]fakeCall{}}
}

type Fake struct {
	calls map[interface{}]fakeCall
}

// Width and height are constrained to zero if negative.
func (i *Fake) NewAcceleratedImage(width, height int) *FakeAcceleratedImage {
	return &FakeAcceleratedImage{
		calls:  i.calls,
		width:  width,
		height: height,
		pixels: make([]Color, width*height),
	}
}

// This method can be used in unit tests for CPU-based functionality
func (i *Fake) NewImageWithFakeAcceleration(width, height int) *Image {
	return New(width, height, i.NewAcceleratedImage(width, height))
}

func (i *Fake) FillWithColor(c Color) AcceleratedCall {
	call := &fillWithColor{color: c}
	i.calls[call] = call
	return call
}

func (i *Fake) NoOp() AcceleratedCall {
	call := &noOp{}
	i.calls[call] = call
	return call
}

type fakeCall interface {
	Run(selection AcceleratedFragmentLocation, image *FakeAcceleratedImage)
}

type fillWithColor struct {
	color Color
}

func (f *fillWithColor) Run(selection AcceleratedFragmentLocation, image *FakeAcceleratedImage) {
	for x := selection.X; x < selection.X+selection.Width; x++ {
		for y := selection.Y; y < selection.Y+selection.Height; y++ {
			index := x + y*image.width
			image.pixels[index] = f.color
		}
	}
}

type noOp struct {
}

func (n noOp) Run(selection AcceleratedFragmentLocation, image *FakeAcceleratedImage) {
}

type FakeAcceleratedImage struct {
	calls  map[interface{}]fakeCall
	pixels []Color
	width  int
	height int
}

func (i *FakeAcceleratedImage) Upload(input AcceleratedFragmentPixels) {
	inputOffset := input.StartingPosition
	location := input.Location
	for y := 0; y < location.Height; y++ {
		for x := 0; x < location.Width; x++ {
			index := y*i.width + x + location.X + location.Y*i.width
			i.pixels[index] = input.Pixels[inputOffset]
			inputOffset += 1
		}
		inputOffset += input.Stride - location.Width
	}
}
func (i *FakeAcceleratedImage) Download(output AcceleratedFragmentPixels) {
	location := output.Location
	outputOffset := output.StartingPosition
	for y := 0; y < location.Height; y++ {
		for x := 0; x < location.Width; x++ {
			index := y*i.width + x + location.X + location.Y*i.width
			output.Pixels[outputOffset] = i.pixels[index]
			outputOffset += 1
		}
		outputOffset += output.Stride - location.Width
	}
}

func (i *FakeAcceleratedImage) Modify(selection AcceleratedFragmentLocation, call AcceleratedCall) {
	fakeCall, ok := i.calls[call]
	if !ok {
		panic("invalid call")
	}
	fakeCall.Run(selection, i)
}
