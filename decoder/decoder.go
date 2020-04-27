// Package decoder provides functionality of static image decoding.
package decoder

import (
	stdimage "image"
	_ "image/gif"
	_ "image/png"
	"os"

	"github.com/jacekolszak/pixiq/image"
)

type ImageFactory interface {
	NewImage(width, height int) *image.Image
}

func New(imageFactory ImageFactory) *Decoder {
	if imageFactory == nil {
		panic("nil imageFactory")
	}
	return &Decoder{imageFactory: imageFactory}
}

type Decoder struct {
	imageFactory ImageFactory
}

// Reader is the equivalent of io.Reader
type Reader interface {
	Read(p []byte) (n int, err error)
}

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
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			target.SetColor(x, y, image.NRGBA(byte(r), byte(g), byte(b), byte(a)))
		}
	}
	return newImage, nil
}

func (d *Decoder) DecodeFile(fileName string) (*image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return d.Decode(file)
}
