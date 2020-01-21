package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/image"
)

var (
	white = image.RGB(255, 255, 255)
	red   = image.RGBA(255, 0, 0, 255)
)

func TestFake_NewAcceleratedImage(t *testing.T) {
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
				img := image.NewFake().NewAcceleratedImage(test.width, test.height)
				assert.NotNil(t, img)
			})
		}
	})
}

func TestFakeAcceleratedImage_Upload(t *testing.T) {
	t.Run("should upload fragment", func(t *testing.T) {
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
				img := image.NewFake().NewAcceleratedImage(test.width, test.height)
				// when
				img.Upload(test.input)
				// then
				output := image.AcceleratedFragmentPixels{
					Location: image.AcceleratedFragmentLocation{Width: test.width, Height: test.height},
					Pixels:   make([]image.Color, len(test.expected)),
					Stride:   test.width,
				}
				img.Download(output)
				assert.Equal(t, test.expected, output.Pixels)
			})
		}
	})

	t.Run("should copy pixels", func(t *testing.T) {
		input := image.AcceleratedFragmentPixels{
			Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
			Pixels:   []image.Color{transparent},
		}
		img := image.NewFake().NewAcceleratedImage(1, 1)
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

func TestFakeAcceleratedImage_Download(t *testing.T) {
	tests := map[string]struct {
		width, height int
		input         []image.Color
		output        image.AcceleratedFragmentPixels
		expected      []image.Color
	}{
		"image 0x0, output 0,0 with size 0x0": {
			output: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{},
				Pixels:   []image.Color{},
			},
			expected: []image.Color{},
		},
		"image 1x1, output 0,0 with size 1x1": {
			width: 1, height: 1,
			input: []image.Color{white},
			output: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				Pixels:   []image.Color{transparent},
			},
			expected: []image.Color{white},
		},
		"image 2x1, output 0,0 with size 1x1": {
			width: 2, height: 1,
			input: []image.Color{white, red},
			output: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				Pixels:   []image.Color{transparent},
			},
			expected: []image.Color{white},
		},
		"image 1x2, output 0,0 with size 1x1": {
			width: 1, height: 2,
			input: []image.Color{white, red},
			output: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				Pixels:   []image.Color{transparent},
			},
			expected: []image.Color{white},
		},
		"image 2x1, output 1,0 with size 1x1": {
			width: 2, height: 1,
			input: []image.Color{white, red},
			output: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{X: 1, Width: 1, Height: 1},
				Pixels:   []image.Color{transparent},
			},
			expected: []image.Color{red},
		},
		"image 1x2, output 0,1 with size 1x1": {
			width: 1, height: 2,
			input: []image.Color{white, red},
			output: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Y: 1, Width: 1, Height: 1},
				Pixels:   []image.Color{transparent},
			},
			expected: []image.Color{red},
		},
		"image 2x2, output 0,1 with size 1x1": {
			width: 2, height: 2,
			input: []image.Color{white, red, white, red},
			output: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Y: 1, Width: 1, Height: 1},
				Pixels:   []image.Color{transparent},
			},
			expected: []image.Color{white},
		},
		"image 1x1, output 0,0 with size 1x1, starting position 1": {
			width: 1, height: 1,
			input: []image.Color{white},
			output: image.AcceleratedFragmentPixels{
				Location:         image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				Pixels:           []image.Color{transparent, transparent},
				StartingPosition: 1,
			},
			expected: []image.Color{transparent, white},
		},
		"image 1x2, output 0,0 with size 1x2, stride 2": {
			width: 1, height: 2,
			input: []image.Color{white, red},
			output: image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: 1, Height: 2},
				Pixels:   []image.Color{transparent, transparent, transparent, transparent},
				Stride:   2,
			},
			expected: []image.Color{white, transparent, red, transparent},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			img := image.NewFake().NewAcceleratedImage(test.width, test.height)
			input := image.AcceleratedFragmentPixels{
				Location: image.AcceleratedFragmentLocation{Width: test.width, Height: test.height},
				Stride:   test.width,
				Pixels:   test.input,
			}
			img.Upload(input)
			// when
			img.Download(test.output)
			// then
			assert.Equal(t, test.expected, test.output.Pixels)
		})
	}
}

func TestFakeAcceleratedImage_Modify(t *testing.T) {
	t.Run("should panic when call has not been created with Fake", func(t *testing.T) {
		fakeImages := image.NewFake()
		img := fakeImages.NewAcceleratedImage(1, 1)
		location := image.AcceleratedFragmentLocation{Width: 1, Height: 1}
		assert.Panics(t, func() {
			img.Modify(location, struct{}{})
		})
	})
	t.Run("AddColor", func(t *testing.T) {
		var (
			colorToAdd = image.RGBA(1, 2, 3, 4)
			color1     = image.RGBA(10, 20, 30, 40)
			modif1     = image.RGBA(11, 22, 33, 44)
			color2     = image.RGBA(50, 60, 70, 80)
			modif2     = image.RGBA(51, 62, 73, 84)
			color3     = image.RGBA(90, 100, 110, 120)
			modif3     = image.RGBA(91, 102, 113, 124)
			color4     = image.RGBA(130, 140, 150, 160)
			modif4     = image.RGBA(131, 142, 153, 164)
		)

		tests := map[string]struct {
			width, height int
			location      image.AcceleratedFragmentLocation
			given         []image.Color
			expected      []image.Color
		}{
			"1x1": {
				width:    1,
				height:   1,
				location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				given:    []image.Color{color1},
				expected: []image.Color{modif1},
			},
			"2x1": {
				width:    2,
				height:   1,
				location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				given:    []image.Color{color1, color2},
				expected: []image.Color{modif1, color2},
			},
			"2x1, location X:1": {
				width:    2,
				height:   1,
				location: image.AcceleratedFragmentLocation{X: 1, Width: 1, Height: 1},
				given:    []image.Color{color1, color2},
				expected: []image.Color{color1, modif2},
			},
			"1x2": {
				width:    1,
				height:   2,
				location: image.AcceleratedFragmentLocation{Width: 1, Height: 1},
				given:    []image.Color{color1, color2},
				expected: []image.Color{modif1, color2},
			},
			"1x2, location Y: 1": {
				width:    1,
				height:   2,
				location: image.AcceleratedFragmentLocation{Y: 1, Width: 1, Height: 1},
				given:    []image.Color{color1, color2},
				expected: []image.Color{color1, modif2},
			},
			"2x2, location Y: 1": {
				width:    2,
				height:   2,
				location: image.AcceleratedFragmentLocation{Y: 1, Width: 1, Height: 1},
				given:    []image.Color{color1, color2, color3, color4},
				expected: []image.Color{color1, color2, modif3, color4},
			},
			"2x1, location Width: 2": {
				width:    2,
				height:   1,
				location: image.AcceleratedFragmentLocation{Width: 2, Height: 1},
				given:    []image.Color{color1, color2},
				expected: []image.Color{modif1, modif2},
			},
			"1x2, location Height: 2": {
				width:    1,
				height:   2,
				location: image.AcceleratedFragmentLocation{Width: 1, Height: 2},
				given:    []image.Color{color1, color2},
				expected: []image.Color{modif1, modif2},
			},
			"2x2, location Width: 2, Height: 2": {
				width:    2,
				height:   2,
				location: image.AcceleratedFragmentLocation{Width: 2, Height: 2},
				given:    []image.Color{color1, color2, color3, color4},
				expected: []image.Color{modif1, modif2, modif3, modif4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				fakeImages := image.NewFake()
				img := fakeImages.NewAcceleratedImage(test.width, test.height)
				img.Upload(image.AcceleratedFragmentPixels{
					Location: image.AcceleratedFragmentLocation{
						X:      0,
						Y:      0,
						Width:  test.width,
						Height: test.height,
					},
					Pixels:           test.given,
					StartingPosition: 0,
					Stride:           test.width,
				})
				// when
				img.Modify(test.location, fakeImages.AddColor(colorToAdd))
				// then
				output := image.AcceleratedFragmentPixels{
					Location: image.AcceleratedFragmentLocation{
						Width:  test.width,
						Height: test.height,
					},
					Stride: test.width,
					Pixels: make([]image.Color, len(test.expected)),
				}
				img.Download(output)
				assert.Equal(t, test.expected, output.Pixels)
			})
		}
	})
	t.Run("RegisterCall", func(t *testing.T) {
		t.Run("should execute custom call", func(t *testing.T) {
			fakeImages := image.NewFake()
			img := fakeImages.NewAcceleratedImage(1, 1)
			location := image.AcceleratedFragmentLocation{Width: 1, Height: 1}
			callMock := &callMock{}
			fakeImages.RegisterCall(callMock)
			// when
			img.Modify(location, callMock)
			// then
			assert.True(t, callMock.executed)
		})
	})
}

type callMock struct {
	executed bool
}

func (f *callMock) Run(selection image.AcceleratedFragmentLocation, image *image.FakeAcceleratedImage) {
	f.executed = true
}
