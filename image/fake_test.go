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
//		selection     image.AcceleratedFragmentLocation
//		output        image.AcceleratedFragmentPixels
//		expected      []image.Color
//	}{
//		"image 1x1, selection 0,0 with size 1x1": {
//			selection: image.AcceleratedFragmentLocation{
//				Width:  1,
//				Height: 1,
//			},
//			output: image.AcceleratedFragmentPixels{
//				Pixels:           make([]image.Color, 1),
//				StartingPosition: 0,
//				Stride:           0,
//			},
//			expected: []image.Color{transparent},
//		},
//		//"image 1x1; selection 0,0 with size 1x1": {
//		//	width: 1, height: 1,
//		//	selection:    image.AcceleratedFragmentLocation{Width: 1, Height: 1},
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
		input         image.AcceleratedFragmentPixels
		expected      []image.Color
	}{
		"image 1x1; fragment 0,0 with size 1x1, white": {
			width: 1, height: 1,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				Pixels:   []image.Color{white},
			},
			expected: []image.Color{white},
		},
		"image 1x1; fragment 0,0 with size 1x1, red": {
			width: 1, height: 1,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				Pixels:   []image.Color{red},
			},
			expected: []image.Color{red},
		},
		"image 1x1; fragment 0,0 with size 0x0": {
			width: 1, height: 1,
			input: image.AcceleratedFragmentPixels{
				Pixels: []image.Color{red},
			},
			expected: []image.Color{transparent},
		},
		"image 0x0; fragment 0,0 with size 0x0": {
			input: image.AcceleratedFragmentPixels{
				Pixels: []image.Color{},
			},
			expected: []image.Color{},
		},
		"image 1x0; fragment 0,0 with size 0x0": {
			width: 1,
			input: image.AcceleratedFragmentPixels{
				Pixels: []image.Color{},
			},
			expected: []image.Color{},
		},
		"image 0x1; fragment 0,0 with size 0x0": {
			height: 1,
			input: image.AcceleratedFragmentPixels{
				Pixels: []image.Color{},
			},
			expected: []image.Color{},
		},
		"image 2x1; fragment 0,0 with size 1x1": {
			width: 2, height: 1,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				Pixels:   []image.Color{red, white},
			},
			expected: []image.Color{red, transparent},
		},
		"image 2x1; fragment 0,0 with size 2x1": {
			width: 2, height: 1,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 2, Height: 1},
				Pixels:   []image.Color{red, white},
			},
			expected: []image.Color{red, white},
		},
		"image 1x2; fragment 0,0 with size 1x1": {
			width: 1, height: 2,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				Pixels:   []image.Color{red, white},
			},
			expected: []image.Color{red, transparent},
		},
		"image 1x2; fragment 0,0 with size 1x2": {
			width: 1, height: 2,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 2},
				Pixels:   []image.Color{red, white},
				Stride:   1,
			},
			expected: []image.Color{red, white},
		},
		"image 2x1; fragment 1,0 with size 1x1": {
			width: 2, height: 1,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{X: 1, Width: 1, Height: 1},
				Pixels:   []image.Color{white},
			},
			expected: []image.Color{transparent, white},
		},
		"image 1x2; fragment 0,1 with size 1x1": {
			width: 1, height: 2,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Y: 1, Width: 1, Height: 1},
				Pixels:   []image.Color{white},
			},
			expected: []image.Color{transparent, white},
		},
		"image 2x2; fragment 0,1 with size 1x1": {
			width: 2, height: 2,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Y: 1, Width: 1, Height: 1},
				Pixels:   []image.Color{white},
			},
			expected: []image.Color{transparent, transparent, white, transparent},
		},
		"image 1x1; fragment 0,0 with size 1x1, starting position 1": {
			width: 1, height: 1,
			input: image.AcceleratedFragmentPixels{
				Location:         image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				StartingPosition: 1,
				Pixels:           []image.Color{white, red},
			},
			expected: []image.Color{red},
		},
		"image 1x2; fragment 0,0 with size 1x2, stride 2": {
			width: 1, height: 2,
			input: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 2},
				Pixels:   []image.Color{white, white, red, white},
				Stride:   2,
			},
			expected: []image.Color{white, red},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			img := image.NewFakeAcceleratedImage(test.width, test.height)
			// when
			img.Upload(test.input)
			// then
			output := image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: test.width, Height: test.height},
				Pixels:   make([]image.Color, len(test.expected)),
			}
			img.Download(output)
			assert.Equal(t, test.expected, output.Pixels)
		})
	}

}

func TestFakeAcceleratedImage_Upload2(t *testing.T) {
	white := image.RGBA(255, 255, 255, 255)

	t.Run("should copy pixels", func(t *testing.T) {
		input := image.AcceleratedFragmentPixels{
			Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
			Pixels:   []image.Color{transparent},
		}
		img := image.NewFakeAcceleratedImage(1, 1)
		img.Upload(input)
		// when
		input.Pixels[0] = white
		// then
		output := image.AcceleratedFragmentPixels{
			Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
			Pixels:   []image.Color{transparent},
		}
		img.Download(output)
		assert.Equal(t, []image.Color{transparent}, output.Pixels)
	})
}
