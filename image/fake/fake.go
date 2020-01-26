package fake

import (
	"errors"

	"github.com/jacekolszak/pixiq/image"
)

// NewAccelerator returns a new container of accelerated images and programs.
// It can be used in unit tests as a replacement for a real implementation
// (such as OpenGL).
func NewAccelerator() *Accelerator {
	return &Accelerator{
		programs: map[image.AcceleratedProgram]*program{},
	}
}

type Accelerator struct {
	programs map[image.AcceleratedProgram]*program
}

func (i *Accelerator) NewImage(imageWidth, imageHeight int) *AcceleratedImage {
	return &AcceleratedImage{
		programs:    i.programs,
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
	}
}

func (i *Accelerator) NewProgram(f func(img *AcceleratedImage, selection image.AcceleratedImageSelection)) image.AcceleratedProgram {
	program := &program{f: f}
	i.programs[program] = program
	return program
}

// TODO Test
// NewAddColorProgram creates a new program adding all color components to each
// pixel in a selection.
func (i *Accelerator) NewAddColorProgram(colorToAdd image.Color) image.AcceleratedProgram {
	return i.NewProgram(func(img *AcceleratedImage, selection image.AcceleratedImageSelection) {
		for y := selection.Y; y < selection.Y+selection.Height; y++ {
			for x := selection.X; x < selection.X+selection.Width; x++ {
				idx := y*img.imageWidth + x
				color := img.pixels[idx]
				var (
					r = color.R() + colorToAdd.R()
					g = color.G() + colorToAdd.G()
					b = color.B() + colorToAdd.B()
					a = color.A() + colorToAdd.A()
				)
				img.pixels[idx] = image.RGBA(r, g, b, a)
			}
		}
	})
}

// AcceleratedImage stores pixel data in RAM and uses CPU solely.
type AcceleratedImage struct {
	pixels      []image.Color
	programs    map[image.AcceleratedProgram]*program
	imageWidth  int
	imageHeight int
}

type drawer struct {
}

func (i *AcceleratedImage) Drawer(program image.AcceleratedProgram, selection image.AcceleratedImageSelection) (image.AcceleratedDrawer, error) {
	prg, ok := i.programs[program]
	if !ok {
		return nil, errors.New("unknown program")
	}
	prg.f(i, selection)
	return &drawer{}, nil
}

func (i *AcceleratedImage) Upload(pixels []image.Color) {
	i.pixels = make([]image.Color, len(pixels))
	// copy pixels to ensure that Upload method has been called
	copy(i.pixels, pixels)
}

func (i *AcceleratedImage) Download(output []image.Color) {
	for j := 0; j < len(output); j++ {
		output[j] = i.pixels[j]
	}
}

type program struct {
	f func(img *AcceleratedImage, selection image.AcceleratedImageSelection)
}
