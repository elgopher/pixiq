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
	return &Decoder{imageFactory: imageFactory}
}

type Decoder struct {
	imageFactory ImageFactory
}

// Reader is the equivalent of io.Reader
type Reader interface {
	Read(p []byte) (n int, err error)
}

func (d *Decoder) Decode(reader Reader) (*DecodedImage, error) {
	img, _, err := stdimage.Decode(reader)
	if err != nil {
		return nil, err
	}
	return &DecodedImage{img: img, imageFactory: d.imageFactory}, nil
}

func (d *Decoder) DecodeFile(fileName string) (*DecodedImage, error) {
	file, _ := os.Open(fileName)
	return d.Decode(file)

}

type DecodedImage struct {
	img          stdimage.Image
	imageFactory ImageFactory
}

func (i *DecodedImage) Width() int {
	return i.img.Bounds().Max.X
}

func (i *DecodedImage) Height() int {
	return i.img.Bounds().Max.Y
}

func (i *DecodedImage) NewImage() *image.Image {
	size := i.img.Bounds().Max
	newImage := i.imageFactory.NewImage(size.X, size.Y)
	i.CopyTo(newImage.WholeImageSelection())
	return newImage
}

func (i *DecodedImage) CopyTo(target image.Selection) {
	size := i.img.Bounds().Max
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			color := i.img.At(x, y)
			r, g, b, a := color.RGBA()
			target.SetColor(x, y, image.NRGBA(byte(r), byte(g), byte(b), byte(a)))
		}
	}
}
