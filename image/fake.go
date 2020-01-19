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
func (i *Fake) FillWithColor(c Color) AcceleratedCall {
	call := &FillWithColor{color: c}
	i.calls[call] = call
	return call
}

type fakeCall interface {
	Run(selection AcceleratedFragmentLocation, image *FakeAcceleratedImage)
}

type FillWithColor struct {
	color Color
}

func (f *FillWithColor) Run(selection AcceleratedFragmentLocation, image *FakeAcceleratedImage) {
	// TODO Implement rest
	image.pixels[0] = f.color
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
