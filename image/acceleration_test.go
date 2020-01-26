package image_test

import (
	"errors"

	"github.com/jacekolszak/pixiq/image"
)

type fakeAcceleratedImage struct {
	pixels   []image.Color
	programs map[image.AcceleratedProgram]*fakeProgram
}

func newFakeAcceleratedImage() *fakeAcceleratedImage {
	return &fakeAcceleratedImage{
		programs: map[image.AcceleratedProgram]*fakeProgram{},
	}
}

type fakeAcceleratedModification struct {
}

func (i *fakeAcceleratedImage) Modify(program image.AcceleratedProgram, selection image.AcceleratedImageSelection) (image.AcceleratedModification, error) {
	prg, ok := i.programs[program]
	if !ok {
		return nil, errors.New("unknown program")
	}
	prg.f(selection)
	return &fakeAcceleratedModification{}, nil
}

func (i *fakeAcceleratedImage) Upload(pixels []image.Color) {
	i.pixels = make([]image.Color, len(pixels))
	// copy pixels to ensure that Upload method has been called
	copy(i.pixels, pixels)
}

func (i *fakeAcceleratedImage) Download(output []image.Color) {
	for j := 0; j < len(output); j++ {
		output[j] = i.pixels[j]
	}
}

func (i *fakeAcceleratedImage) NewProgram(f func(selection image.AcceleratedImageSelection)) image.AcceleratedProgram {
	program := &fakeProgram{f: f}
	i.programs[program] = program
	return program
}

type fakeProgram struct {
	f func(selection image.AcceleratedImageSelection)
}
