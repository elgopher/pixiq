package fake

import (
	"errors"

	"github.com/jacekolszak/pixiq/image"
)

// NewAccelerator returns a new instance of Accelerator.
func NewAccelerator() *Accelerator {
	return &Accelerator{
		programs: map[image.AcceleratedProgram]*Program{},
	}
}

// Accelerator is a container of accelerated images and programs.
// It can be used in unit tests as a replacement for a real implementation
// (such as OpenGL).
type Accelerator struct {
	programs map[image.AcceleratedProgram]*Program
}

type Primitive struct {
	image.Primitive
	drawn bool
}

func (p *Primitive) Drawn() bool {
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
// a drawer.
func (i *Accelerator) NewProgram() *Program {
	program := &Program{}
	i.programs[program] = program
	return program
}

// AcceleratedImage stores pixel data in RAM and uses CPU solely.
type AcceleratedImage struct {
	pixels      []image.Color
	programs    map[image.AcceleratedProgram]*Program
	imageWidth  int
	imageHeight int
}

type drawer struct {
	selections     map[string]image.AcceleratedImageSelection
	targetLocation image.AcceleratedImageLocation
	targetImage    *AcceleratedImage
}

func (d *drawer) Draw(primitive image.Primitive, params ...interface{}) error {
	fakePrimitive, ok := primitive.(*Primitive)
	if !ok {
		return errors.New("primitive cannot be drawn")
	}
	fakePrimitive.drawn = true
	return nil
}

func (d *drawer) SetSelection(name string, selection image.AcceleratedImageSelection) {
	d.selections[name] = selection
}

func (i *AcceleratedImage) Modify(program image.AcceleratedProgram, location image.AcceleratedImageLocation, procedure func(drawer image.AcceleratedDrawer)) error {
	prg, ok := i.programs[program]
	if !ok {
		return errors.New("unknown program")
	}
	drawer := &drawer{
		selections:     map[string]image.AcceleratedImageSelection{},
		targetLocation: location,
		targetImage:    i,
	}
	prg.executed = true
	prg.drawer = drawer
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

type Program struct {
	drawer   *drawer
	executed bool
}

func (p *Program) Executed() bool {
	return p.executed
}

func (p *Program) TargetLocation() image.AcceleratedImageLocation {
	return p.drawer.targetLocation
}

func (p *Program) TargetImage() image.AcceleratedImage {
	return p.drawer.targetImage
}
