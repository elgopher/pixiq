package pixiq

// NewAcceleratedImage creates an accelerated image. This can be a texture on a video card or something totally different.
type NewAcceleratedImage func(width, height int) AcceleratedImage

// AcceleratedImage is an image processed externally (outside the CPU)
type AcceleratedImage interface {
	Upload(pixels []Color)
}
