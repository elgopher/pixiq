package decoder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
