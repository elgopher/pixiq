package pixiq

type NewAcceleratedImage func(width, height int) AcceleratedImage

type AcceleratedImage interface {
	Upload(pixels []Color)
}
