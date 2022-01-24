// Package decoder provides functionality of decoding compressed images such as
// PNG and GIF.
package decoder

import (
	stdimage "image"
	_ "image/gif" // Register GIF decoder
	_ "image/png" // Register PNG decoder
	"os"

	"github.com/elgopher/pixiq/image"
)

// ImageFactory creates a new image with given dimensions.
//
// *glfw.OpenGL instance can be used as an ImageFactory implementation.
type ImageFactory interface {
	NewImage(width, height int) *image.Image
}

// New creates a Decoder instance which can be used many times for image decoding.
func New(imageFactory ImageFactory) *Decoder {
	if imageFactory == nil {
		panic("nil imageFactory")
	}
	return &Decoder{imageFactory: imageFactory}
}

// Decoder decodes compressed images, such as PNGs and GIFs
type Decoder struct {
	imageFactory ImageFactory
}

// Reader is the equivalent of io.Reader
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Decode decodes compressed image such as PNG or GIF and creates a new *image.Image
// object filled with colors from decompressed image.
func (d *Decoder) Decode(reader Reader) (*image.Image, error) {
	if reader == nil {
		panic("nil reader")
	}
	img, _, err := stdimage.Decode(reader)
	if err != nil {
		return nil, err
	}
	size := img.Bounds().Max
	newImage := d.imageFactory.NewImage(size.X, size.Y)
	target := newImage.WholeImageSelection()
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			color := image.RGBA(byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
			target.SetColor(x, y, color)
		}
	}
	return newImage, nil
}

// DecodeFile decodes compressed file such as PNG or GIF and creates a new *image.Image
// object filled with colors from decompressed file.
func (d *Decoder) DecodeFile(fileName string) (*image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return d.Decode(file)
}
