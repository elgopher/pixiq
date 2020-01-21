package image

func NewFakeImages() *FakeImages {
	return &FakeImages{calls: map[interface{}]FakeCall{}}
}

type FakeImages struct {
	calls map[interface{}]FakeCall
}

// Width and height are constrained to zero if negative.
func (i *FakeImages) NewAcceleratedImage(width, height int) *FakeAcceleratedImage {
	return &FakeAcceleratedImage{
		calls:  i.calls,
		width:  width,
		height: height,
		pixels: make([]Color, width*height),
	}
}

// This method can be used in unit tests for CPU-based functionality
func (i *FakeImages) NewImageWithFakeAcceleration(width, height int) *Image {
	return New(width, height, i.NewAcceleratedImage(width, height))
}

func (i *FakeImages) RegisterCall(call FakeCall) {
	if call == nil {
		panic("nil call")
	}
	i.calls[call] = call
}

type FakeCall interface {
	Run(selection AcceleratedFragmentLocation, image *FakeAcceleratedImage)
}

func (i *FakeImages) AddColor(c Color) AcceleratedCall {
	call := &fakeAddColor{color: c}
	i.calls[call] = call
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

type FakeAcceleratedImage struct {
	calls  map[interface{}]FakeCall
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
			inputOffset++
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
			outputOffset++
		}
		outputOffset += output.Stride - location.Width
	}
}

func (i *FakeAcceleratedImage) Modify(selection AcceleratedFragmentLocation, call AcceleratedCall) {
	if call == nil {
		panic("nil call")
	}
	fakeCall, ok := i.calls[call]
	if !ok {
		panic("invalid call")
	}
	fakeCall.Run(selection, i)
}
