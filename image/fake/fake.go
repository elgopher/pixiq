package fake

import (
	"errors"

	"github.com/jacekolszak/pixiq/image"
)

// NewAccelerator returns a new instance of Accelerator.
func NewAccelerator() *Accelerator {
	return &Accelerator{
		programs: map[image.AcceleratedProgram]*program{},
	}
}

// Accelerator is a container of accelerated images and programs.
// It can be used in unit tests as a replacement for a real implementation
// (such as OpenGL).
type Accelerator struct {
	programs map[image.AcceleratedProgram]*program
}

type Primtive struct {
	image.Primitive
	drawn bool
}

func (p *Primtive) Drawn() bool {
	return p.drawn
}

// NewImage returns a new instance of *AcceleratedImage
func (i *Accelerator) NewImage(imageWidth, imageHeight int) *AcceleratedImage {
	img := &AcceleratedImage{
		programs:    i.programs,
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
	}
	return img
}

// NewProgram returns a new instance of program which can be used to create
// a Drawer.
func (i *Accelerator) NewProgram() image.AcceleratedProgram {
	program := &program{}
	i.programs[program] = program
	return program
}

// AcceleratedImage stores pixel data in RAM and uses CPU solely.
type AcceleratedImage struct {
	pixels      []image.Color
	programs    map[image.AcceleratedProgram]*program
	imageWidth  int
	imageHeight int
}

type Drawer struct {
	selections  map[string]image.AcceleratedImageSelection
	program     *program
	location    image.AcceleratedImageLocation
	outputImage *AcceleratedImage
}

func (d *Drawer) Draw(primitive image.Primitive, params ...interface{}) error {
	fakePrimitive, ok := primitive.(*Primtive)
	if !ok {
		return errors.New("primitive cannot be drawn")
	}
	fakePrimitive.drawn = true
	return nil
}

func (d *Drawer) SetSelection(name string, selection image.AcceleratedImageSelection) {
	d.selections[name] = selection
}

// Drawer returns an AcceleratedDrawer for the program created by the same Accelerator
// as this image.
func (i *AcceleratedImage) Modify(program image.AcceleratedProgram, location image.AcceleratedImageLocation, procedure func(drawer image.AcceleratedDrawer)) error {
	prg, ok := i.programs[program]
	if !ok {
		return errors.New("unknown program")
	}
	drawer := &Drawer{
		selections:  map[string]image.AcceleratedImageSelection{},
		program:     prg,
		location:    location,
		outputImage: i,
	}
	procedure(drawer)
	return nil
}

// Upload send pixels to a container in RAM
func (i *AcceleratedImage) Upload(pixels []image.Color) {
	i.pixels = make([]image.Color, len(pixels))
	// copy pixels to ensure that Upload method has been called
	copy(i.pixels, pixels)
}

// Download fills output slice with image colors
func (i *AcceleratedImage) Download(output []image.Color) {
	for j := 0; j < len(output); j++ {
		output[j] = i.pixels[j]
	}
}

type program struct{}
