package decoder_test

import (
	"bytes"
	stdimage "image"
	"image/color"
	"image/gif"
	"image/png"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			"png 1x2":              png1x2(),
			"png 2x1":              png2x1(),
			"gif 1x2":              gif1x2(),
			"png semi-transparent": pngSemiTransparent(),
			"png 64-bit":           png64bit(),
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
						expectedColor := test.expectedColors[y][x]
						actualColor := selection.Color(x, y)
						assertColor(t, expectedColor, actualColor)
					}
				}
			})
		}

	})
}

func TestDecoder_DecodeFile(t *testing.T) {
	t.Run("should return error", func(t *testing.T) {
		filenames := []string{"", "not-existing-file"}
		for _, filename := range filenames {
			t.Run(filename, func(t *testing.T) {
				imageDecoder := decoder.New(fakeImageFactory{})
				// when
				img, err := imageDecoder.DecodeFile(filename)
				assert.Error(t, err)
				assert.Nil(t, img)
			})
		}
	})
	t.Run("should decode file", func(t *testing.T) {
		imageDecoder := decoder.New(fakeImageFactory{})
		// when
		img, err := imageDecoder.DecodeFile(pngFilename(t))
		require.NoError(t, err)
		assert.NotNil(t, img)
	})
}

func pngFilename(t *testing.T) string {
	file, err := ioutil.TempFile("", "TestDecoder_DecodeFile")
	require.NoError(t, err)
	defer file.Close()
	pngImage := stdimage.NewNRGBA(stdimage.Rect(0, 0, 1, 1))
	buffer := bytes.Buffer{}
	_ = png.Encode(&buffer, pngImage)
	_, err = file.Write(buffer.Bytes())
	require.NoError(t, err)
	return file.Name()
}

func assertColor(t *testing.T, expectedColor image.Color, actualColor image.Color) {
	r, g, b, a := expectedColor.RGBA()
	ra, ga, ba, aa := actualColor.RGBA()
	delta := 1.0
	assert.InDelta(t, r, ra, delta, "red components differ, expected: %v, actual %v", expectedColor, actualColor)
	assert.InDelta(t, g, ga, delta, "green components differ, expected: %v, actual %v", expectedColor, actualColor)
	assert.InDelta(t, b, ba, delta, "blue components differ, expected: %v, actual %v", expectedColor, actualColor)
	assert.InDelta(t, a, aa, delta, "alpha components differ, expected: %v, actual %v", expectedColor, actualColor)
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
	return testCase{
		data: buffer.Bytes(),
		expectedColors: [][]image.Color{
			{colornames.Black},
			{colornames.White},
		}}
}
func png2x1() testCase {
	pngImage := stdimage.NewNRGBA(stdimage.Rect(0, 0, 2, 1))
	pngImage.Set(0, 0, stdimage.Black)
	pngImage.Set(1, 0, stdimage.White)
	buffer := bytes.Buffer{}
	_ = png.Encode(&buffer, pngImage)
	return testCase{
		data: buffer.Bytes(),
		expectedColors: [][]image.Color{
			{colornames.Black, colornames.White},
		},
	}
}

func gif1x2() testCase {
	gifImage := stdimage.NewNRGBA(stdimage.Rect(0, 0, 1, 2))
	gifImage.Set(0, 0, stdimage.Black)
	gifImage.Set(0, 1, stdimage.White)
	buffer := bytes.Buffer{}
	_ = gif.Encode(&buffer, gifImage, &gif.Options{
		NumColors: 256,
	})
	return testCase{
		data: buffer.Bytes(),
		expectedColors: [][]image.Color{
			{colornames.Black},
			{colornames.White},
		},
	}
}

func pngSemiTransparent() testCase {
	pngImage := stdimage.NewRGBA(stdimage.Rect(0, 0, 1, 1))
	pngImage.Set(0, 0, color.RGBA{R: 50, G: 100, B: 150, A: 200})
	buffer := bytes.Buffer{}
	_ = png.Encode(&buffer, pngImage)
	return testCase{
		data: buffer.Bytes(),
		expectedColors: [][]image.Color{
			{image.RGBA(50, 100, 150, 200)},
		}}
}

func png64bit() testCase {
	pngImage := stdimage.NewRGBA64(stdimage.Rect(0, 0, 1, 1))
	pngImage.Set(0, 0, color.RGBA64{R: 5000, G: 10000, B: 15000, A: 20000})
	buffer := bytes.Buffer{}
	_ = png.Encode(&buffer, pngImage)
	return testCase{
		data: buffer.Bytes(),
		expectedColors: [][]image.Color{
			{image.RGBA(20, 39, 59, 78)},
		},
	}
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
