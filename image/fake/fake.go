package fake

import (
	"errors"

	"github.com/jacekolszak/pixiq/image"
)

type Primitive struct {
	image.Primitive
	drawn      bool
	params     []interface{}
	selections []Selection
}

type Selection struct {
	name     string
	location image.AcceleratedImageLocation
	image    *AcceleratedImage
}

func (s *Selection) Location() image.AcceleratedImageLocation {
	return s.location
}

func (s *Selection) AcceleratedImage() *AcceleratedImage {
	return s.image
}

func (s *Selection) Name() string {
	return s.name
}

func (p *Primitive) Drawn() bool {
	return p.drawn
}

func (p *Primitive) ParamsPassed() []interface{} {
	return p.params
}

func (p *Primitive) SelectionsPassed() []Selection {
	return p.selections
}

// NewAcceleratedImage returns a new instance of *AcceleratedImage
func NewAcceleratedImage(imageWidth, imageHeight int) *AcceleratedImage {
	img := &AcceleratedImage{
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
		pixels:      make([]image.Color, imageWidth*imageHeight),
	}
	return img
}

func NewProgram() *Program {
	return &Program{}
}

// AcceleratedImage stores pixel data in RAM and uses CPU solely.
type AcceleratedImage struct {
	pixels      []image.Color
	imageWidth  int
	imageHeight int
}

type drawer struct {
	selections     map[string]image.AcceleratedImageSelection
	targetLocation image.AcceleratedImageLocation
	targetImage    *AcceleratedImage
}

func (d *drawer) Draw(primitive image.Primitive, params []interface{}) error {
	fakePrimitive, ok := primitive.(*Primitive)
	if !ok {
		return errors.New("primitive cannot be drawn")
	}
	fakePrimitive.drawn = true
	fakePrimitive.params = params
	for name, selection := range d.selections {
		fakeAcceleratedImage, ok := selection.AcceleratedImage.(*AcceleratedImage)
		if !ok {
			return errors.New("selection's AcceleratedImage is not *fake.AcceleratedImage")
		}
		// TODO Take a snapshot of image and export in Selection
		fakePrimitive.selections = append(fakePrimitive.selections, Selection{
			name:     name,
			location: selection.AcceleratedImageLocation,
			image:    fakeAcceleratedImage,
		})
	}
	return nil
}

func (d *drawer) SetSelection(name string, selection image.AcceleratedImageSelection) {
	d.selections[name] = selection
}

func (i *AcceleratedImage) Modify(program image.AcceleratedProgram, location image.AcceleratedImageLocation, procedure func(drawer image.AcceleratedDrawer)) error {
	fakeProgram, ok := program.(*Program)
	if !ok {
		return errors.New("cannot execute a program which is not a *fake.Program")
	}
	drawer := &drawer{
		selections:     map[string]image.AcceleratedImageSelection{},
		targetLocation: location,
		targetImage:    i,
	}
	fakeProgram.executed = true
	fakeProgram.drawer = drawer
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
