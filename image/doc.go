// Package image provides hardware-supported Image primitive which can be
// manipulated in real-time using Pixel-Perfect API.
//
// Instance of the Image should be created either by directly using a New function
// or using an external factory function such as opengl.OpenGL.NewImage:
//
//	   image := images.New(2, 2, acceleratedImage)
//
// Image can be manipulated using Selection:
//
//     wholeSelection := image.WholeImageSelection()
//     wholeSelection.SetColor(0, 1, colornames.White)
//
package image
