// Package pixiq allows you to manipulate images in real-time using Pixel-Perfect API.
//
// Before Image can be created a factory should be created:
//
//	   images := pixiq.NewImages(acceleratedImages) // i.e. gl.AcceleratedImages()
//
// Then multiple images can be created:
//
//	   image := images.New(2, 2)
//
// And pixel set using Selection:
//
//     wholeSelection := image.WholeImageSelection()
// 	   wholeSelection.SetColor(0, 1, colornames.White)
//
package pixiq
