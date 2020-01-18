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
	Selection        AcceleratedFragment
	Call             AcceleratedCall
	PixelsDuringCall []Color
}

func (i *FakeAcceleratedImage) Upload(selection AcceleratedFragment, input PixelSlice) {
	if i.pixels == nil {
		i.pixels = make([]Color, i.width*i.height)
	}
	inputOffset := input.StartingPosition
	for y := 0; y < selection.Height; y++ {
		for x := 0; x < selection.Width; x++ {
			index := y*i.width + x + selection.X + selection.Y*i.width
			i.pixels[index] = input.Pixels[inputOffset]
			inputOffset += 1
		}
		inputOffset += input.Stride - selection.Width
	}
}
func (i *FakeAcceleratedImage) Download(selection AcceleratedFragment, pixels PixelSlice) {
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

func (i *FakeAcceleratedImage) Modify(selection AcceleratedFragment, call AcceleratedCall) {
	pixelsCopy := make([]Color, len(i.pixels))
	copy(pixelsCopy, i.pixels)
	modifyCall := ModifyCall{
		Selection:        selection,
		Call:             call,
		PixelsDuringCall: pixelsCopy,
	}
	i.ModifyCalls = append(i.ModifyCalls, modifyCall)
}
