package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/image"
)

func TestNewFakeAcceleratedImage(t *testing.T) {
	t.Run("should create FakeAcceleratedImage for testing", func(t *testing.T) {
		tests := map[string]struct {
			width, height int
		}{
			"0 x 0": {
				width:  0,
				height: 0,
			},
			"-1 x 0": {
				width:  -1,
				height: 0,
			},
			"0 x -1": {
				width:  0,
				height: -1,
			},
			"1 x 1": {
				width:  0,
				height: 0,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				img := image.NewFakeAcceleratedImage(test.width, test.height)
				assert.NotNil(t, img)
			})
		}
	})
}

//
//func TestFakeAcceleratedImage_Download(t *testing.T) {
//	//white := image.RGBA(255, 255, 255, 255)
//	tests := map[string]struct {
//		width, height int
//		selection     image.AcceleratedSelection
//		output        image.PixelSlice
//		expected      []image.Color
//	}{
//		"image 1x1, selection 0,0 with size 1x1": {
//			selection: image.AcceleratedSelection{
//				Width:  1,
//				Height: 1,
//			},
//			output: image.PixelSlice{
//				Pixels:           make([]image.Color, 1),
//				StartingPosition: 0,
//				Stride:           0,
//			},
//			expected: []image.Color{transparent},
//		},
//		//"image 1x1; selection 0,0 with size 1x1": {
//		//	width: 1, height: 1,
//		//	selection:    image.AcceleratedSelection{Width: 1, Height: 1},
//		//	outputPixels: []image.Color{white},
//		//	expected:     []image.Color{white},
//		//},
//	}
//	for name, test := range tests {
//		t.Run(name, func(t *testing.T) {
//			img := image.NewFakeAcceleratedImage(test.width, test.height)
//			// when
//			img.Download2(test.selection, test.output)
//			// then
//			assert.Equal(t, test.expected, test.output.Pixels)
//		})
//	}
//}

func TestFakeAcceleratedImage_Upload(t *testing.T) {
	white := image.RGBA(255, 255, 255, 255)
	red := image.RGBA(255, 0, 0, 255)
	tests := map[string]struct {
		width, height int
		selection     image.AcceleratedSelection
		input         image.PixelSlice
		expected      []image.Color
	}{
		"image 1x1; selection 0,0 with size 1x1, white": {
			width: 1, height: 1,
			selection: image.AcceleratedSelection{Width: 1, Height: 1},
			input: image.PixelSlice{
				Pixels: []image.Color{white},
			},
			expected: []image.Color{white},
		},
		"image 1x1; selection 0,0 with size 1x1, red": {
			width: 1, height: 1,
			selection: image.AcceleratedSelection{Width: 1, Height: 1},
			input: image.PixelSlice{
				Pixels: []image.Color{red},
			},
			expected: []image.Color{red},
		},
		"image 1x1; selection 0,0 with size 0x0": {
			width: 1, height: 1,
			input: image.PixelSlice{
				Pixels: []image.Color{red},
			},
			expected: []image.Color{transparent},
		},
		"image 0x0; selection 0,0 with size 0x0": {
			input: image.PixelSlice{
				Pixels: []image.Color{},
			},
			expected: []image.Color{},
		},
		"image 1x0; selection 0,0 with size 0x0": {
			width: 1,
			input: image.PixelSlice{
				Pixels: []image.Color{},
			},
			expected: []image.Color{},
		},
		"image 0x1; selection 0,0 with size 0x0": {
			height: 1,
			input: image.PixelSlice{
				Pixels: []image.Color{},
			},
			expected: []image.Color{},
		},
		"image 2x1; selection 0,0 with size 1x1": {
			width: 2, height: 1,
			selection: image.AcceleratedSelection{Width: 1, Height: 1},
			input: image.PixelSlice{
				Pixels: []image.Color{red, white},
			},
			expected: []image.Color{red, transparent},
		},
		"image 2x1; selection 0,0 with size 2x1": {
			width: 2, height: 1,
			selection: image.AcceleratedSelection{Width: 2, Height: 1},
			input: image.PixelSlice{
				Pixels: []image.Color{red, white},
			},
			expected: []image.Color{red, white},
		},
		"image 1x2; selection 0,0 with size 1x1": {
			width: 1, height: 2,
			selection: image.AcceleratedSelection{Width: 1, Height: 1},
			input: image.PixelSlice{
				Pixels: []image.Color{red, white},
			},
			expected: []image.Color{red, transparent},
		},
		"image 1x2; selection 0,0 with size 1x2": {
			width: 1, height: 2,
			selection: image.AcceleratedSelection{Width: 1, Height: 2},
			input: image.PixelSlice{
				Pixels: []image.Color{red, white},
			},
			expected: []image.Color{red, white},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			img := image.NewFakeAcceleratedImage(test.width, test.height)
			// when
			img.Upload(test.selection, test.input)
			// then
			whole := image.AcceleratedSelection{Width: test.width, Height: test.height}
			output := image.PixelSlice{
				Pixels: make([]image.Color, len(test.expected)),
			}
			img.Download(whole, output)
			assert.Equal(t, test.expected, output.Pixels)
		})
	}

}

func TestFakeAcceleratedImage_Upload2(t *testing.T) {
	white := image.RGBA(255, 255, 255, 255)

	t.Run("should copy pixels", func(t *testing.T) {
		input := image.PixelSlice{
			Pixels: []image.Color{transparent},
		}
		img := image.NewFakeAcceleratedImage(1, 1)
		selection := image.AcceleratedSelection{Width: 1, Height: 1}
		img.Upload(selection, input)
		// when
		input.Pixels[0] = white
		// then
		output := image.PixelSlice{
			Pixels: []image.Color{transparent},
		}
		img.Download(selection, output)
		assert.Equal(t, []image.Color{transparent}, output.Pixels)
	})
}
