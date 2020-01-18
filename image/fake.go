package image

// Width and height are constrained to zero if negative.
func NewFakeAcceleratedImage(width, height int) *FakeAcceleratedImage {
	return &FakeAcceleratedImage{width: width, height: height}
}

type FakeAcceleratedImage struct {
	pixels      []Color
	width       int
	height      int
	ModifyCalls []ModifyCall
}

type ModifyCall struct {
	Selection        AcceleratedFragmentLocation
	Call             AcceleratedCall
	PixelsDuringCall []Color
}

func (i *FakeAcceleratedImage) Upload(input AcceleratedFragmentPixels) {
	if i.pixels == nil {
		i.pixels = make([]Color, i.width*i.height)
	}
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
func (i *FakeAcceleratedImage) Download(pixels AcceleratedFragmentPixels) {
	if i.width == 0 || i.height == 0 {
		return
	}
	for y := 0; y < i.height; y++ {
		for x := 0; x < i.width; x++ {
			idx := y*i.width + x
			pixels.Pixels[idx] = i.pixels[idx]
		}
	}
}

func (i *FakeAcceleratedImage) Modify(selection AcceleratedFragmentLocation, call AcceleratedCall) {
	pixelsCopy := make([]Color, len(i.pixels))
	copy(pixelsCopy, i.pixels)
	modifyCall := ModifyCall{
		Selection:        selection,
		Call:             call,
		PixelsDuringCall: pixelsCopy,
	}
	i.ModifyCalls = append(i.ModifyCalls, modifyCall)
}
