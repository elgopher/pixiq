package gl

import (
	"fmt"
	"github.com/jacekolszak/pixiq/image"
)

// NewAcceleratedImage returns an OpenGL-accelerated implementation of image.AcceleratedImage
// Will panic if width or height are negative or higher than MAX_TEXTURE_SIZE
func (c *Context) NewAcceleratedImage(width, height int) *AcceleratedImage {
	if width < 0 {
		panic("negative width")
	}
	if height < 0 {
		panic("negative height")
	}
	if width > c.capabilities.maxTextureSize {
		panic(fmt.Sprintf("width higher than MAX_TEXTURE_SIZE (%d pixels)", c.capabilities.maxTextureSize))
	}
	if height > c.capabilities.maxTextureSize {
		panic(fmt.Sprintf("height higher than MAX_TEXTURE_SIZE (%d pixels)", c.capabilities.maxTextureSize))
	}
	// FIXME resize image (internally) if OpenGL does support only a power-of-two dimensions.
	var id uint32
	var frameBufferID uint32
	c.api.GenTextures(1, &id)
	c.api.BindTexture(texture2D, id)
	c.api.TexImage2D(
		texture2D,
		0,
		rgba,
		int32(width),
		int32(height),
		0,
		rgba,
		unsignedByte,
		Ptr(nil),
	)
	c.api.TexParameteri(texture2D, textureMinFilter, nearest)
	c.api.TexParameteri(texture2D, textureMagFilter, nearest)

	c.api.GenFramebuffers(1, &frameBufferID)
	c.api.BindFramebuffer(framebuffer, frameBufferID)
	c.api.FramebufferTexture2D(framebuffer, colorAttachment0, texture2D, id, 0)
	img := &AcceleratedImage{
		textureID:     id,
		frameBufferID: frameBufferID,
		width:         width,
		height:        height,
		api:           c.api,
	}
	c.allImages[img] = img
	return img
}

// AcceleratedImage is an image.AcceleratedImage implementation storing pixels
// on a video card VRAM.
type AcceleratedImage struct {
	textureID     uint32
	frameBufferID uint32
	width, height int
	api           API
}

// TextureID returns texture identifier (aka name)
func (t *AcceleratedImage) TextureID() uint32 {
	return t.textureID
}

// Upload send pixels to video card
func (t *AcceleratedImage) Upload(pixels []image.Color) {
	if len(pixels) == 0 {
		return
	}
	t.api.BindTexture(texture2D, t.textureID)
	t.api.TexSubImage2D(
		texture2D,
		0,
		int32(0),
		int32(0),
		int32(t.width),
		int32(t.height),
		rgba,
		unsignedByte,
		Ptr(pixels),
	)
}

// Download gets pixels pixels from video card
func (t *AcceleratedImage) Download(output []image.Color) {
	if len(output) == 0 {
		return
	}
	t.api.BindTexture(texture2D, t.textureID)
	t.api.GetTexImage(
		texture2D,
		0,
		rgba,
		unsignedByte,
		Ptr(output),
	)
}
