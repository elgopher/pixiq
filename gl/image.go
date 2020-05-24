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
		c.api.Ptr(nil),
	)
	c.api.TexParameteri(texture2D, textureMinFilter, nearest)
	c.api.TexParameteri(texture2D, textureMagFilter, nearest)
	c.api.TexParameteri(texture2D, textureWrapS, clampToBorder)
	c.api.TexParameteri(texture2D, textureWrapT, clampToBorder)

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
	clearWithTransparentColor := c.NewClearCommand()
	clearWithTransparentColor.Run(img.wholeSelection(), []image.AcceleratedImageSelection{})
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

func (i *AcceleratedImage) wholeSelection() image.AcceleratedImageSelection {
	return image.AcceleratedImageSelection{
		Location: image.AcceleratedImageLocation{Width: i.width, Height: i.height},
		Image:    i,
	}
}

// TextureID returns texture identifier (aka name)
func (i *AcceleratedImage) TextureID() uint32 {
	return i.textureID
}

// Upload send pixels to video card
func (i *AcceleratedImage) Upload(pixels []image.Color) {
	if len(pixels) == 0 {
		return
	}
	i.api.BindTexture(texture2D, i.textureID)
	i.api.TexSubImage2D(
		texture2D,
		0,
		int32(0),
		int32(0),
		int32(i.width),
		int32(i.height),
		rgba,
		unsignedByte,
		i.api.Ptr(pixels),
	)
}

// Download gets pixels pixels from video card
func (i *AcceleratedImage) Download(output []image.Color) {
	if len(output) == 0 {
		return
	}
	i.api.BindTexture(texture2D, i.textureID)
	i.api.GetTexImage(
		texture2D,
		0,
		rgba,
		unsignedByte,
		i.api.Ptr(output),
	)
}

// Width returns the number of pixels in a row.
func (i *AcceleratedImage) Width() int {
	return i.width
}

// Height returns the number of pixels in a column.
func (i *AcceleratedImage) Height() int {
	return i.height
}

// Delete cleans resources, usually pixels stored in external memory (such as VRAM).
// After AcceleratedImage is deleted it cannot be used anymore. Subsequent calls
// will generate OpenGL error which can be returned by executing Context.Error()
func (i *AcceleratedImage) Delete() {
	i.api.DeleteTextures(1, &i.textureID)
	i.api.DeleteBuffers(1, &i.frameBufferID)
}
