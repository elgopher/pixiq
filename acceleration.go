package pixiq

// AcceleratedImages is a container of accelerated images.
type AcceleratedImages interface {
	// New creates an accelerated image. This can be a texture on a video card or something totally different.
	New(width, height int) AcceleratedImage
}

// AcceleratedImage is an image processed externally (outside the CPU).
type AcceleratedImage interface {
	// Upload send pixels colors sorted by coordinates. First all pixels are sent for y=0, from left to right.
	Upload(pixels []Color)
}
