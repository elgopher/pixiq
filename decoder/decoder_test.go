package decoder_test

import (
	"bytes"
	stdimage "image"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/decoder"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
)

func TestNew(t *testing.T) {
	t.Run("should panic for nil ImageFactory", func(t *testing.T) {
		assert.Panics(t, func() {
			decoder.New(nil)
		})
	})
	t.Run("should create Decoder", func(t *testing.T) {
		imageDecoder := decoder.New(fakeImageFactory{})
		assert.NotNil(t, imageDecoder)
	})
}

func TestDecoder_Decode(t *testing.T) {
	t.Run("should panic for nil reader", func(t *testing.T) {
		imageDecoder := decoder.New(fakeImageFactory{})
		assert.Panics(t, func() {
			_, _ = imageDecoder.Decode(nil)
		})
	})
	t.Run("should return error when reader returned error", func(t *testing.T) {
		imageDecoder := decoder.New(fakeImageFactory{})
		img, err := imageDecoder.Decode(&erroneousReader{})
		assert.Nil(t, img)
		assert.Error(t, err)
	})
	t.Run("should return error when reader has invalid format", func(t *testing.T) {
		imageDecoder := decoder.New(fakeImageFactory{})
		img, err := imageDecoder.Decode(&invalidFormatReader{})
		assert.Nil(t, img)
		assert.Error(t, err)
	})
	t.Run("should decode image", func(t *testing.T) {
		tests := map[string]testCase{
			"png 1x2": png1x2(),
			"png 2x1": png2x1(),
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				imageDecoder := decoder.New(fakeImageFactory{})
				// when
				img, err := imageDecoder.Decode(bytes.NewReader(test.data))
				// then
				assert.NotNil(t, img)
				assert.NoError(t, err)
				// and
				width := len(test.expectedColors[0])
				assert.Equal(t, img.Width(), width)
				height := len(test.expectedColors)
				assert.Equal(t, img.Height(), height)
				// and
				selection := img.WholeImageSelection()
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						assert.Equal(t, test.expectedColors[y][x], selection.Color(x, y))
					}
				}
			})
		}

	})
}

type testCase struct {
	data           []byte
	expectedColors [][]image.Color
}

func png1x2() testCase {
	pngImage := stdimage.NewNRGBA(stdimage.Rect(0, 0, 1, 2))
	pngImage.Set(0, 0, stdimage.Black)
	pngImage.Set(0, 1, stdimage.White)
	buffer := bytes.Buffer{}
	_ = png.Encode(&buffer, pngImage)
	return testCase{buffer.Bytes(), [][]image.Color{{colornames.Black}, {colornames.White}}}
}
func png2x1() testCase {
	pngImage := stdimage.NewNRGBA(stdimage.Rect(0, 0, 2, 1))
	pngImage.Set(0, 0, stdimage.Black)
	pngImage.Set(1, 0, stdimage.White)
	buffer := bytes.Buffer{}
	_ = png.Encode(&buffer, pngImage)
	return testCase{buffer.Bytes(), [][]image.Color{{colornames.Black, colornames.White}}}
}

type fakeImageFactory struct {
}

func (i fakeImageFactory) NewImage(width, height int) *image.Image {
	return image.New(fake.NewAcceleratedImage(width, height))
}

type erroneousReader struct {
	error error
}

func (e *erroneousReader) Read(p []byte) (n int, err error) {
	return 0, e.error
}

type invalidFormatReader struct {
}

func (i invalidFormatReader) Read(p []byte) (n int, err error) {
	return 0, nil
}
