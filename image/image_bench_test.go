package image_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/image"
)

func BenchmarkSelection_SetColor(b *testing.B) {
	var (
		color     = image.RGBA(10, 20, 30, 40)
		img       = image.New(acceleratedImageStub{width: 1920, height: 1080})
		selection = img.WholeImageSelection()
		height    = selection.Height()
		width     = selection.Width()
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				selection.SetColor(x, y, color)
			}
		}
	}
}

func BenchmarkSelection_Color(b *testing.B) {
	var (
		img       = image.New(acceleratedImageStub{width: 1920, height: 1080})
		selection = img.WholeImageSelection()
		height    = selection.Height()
		width     = selection.Width()
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				selection.Color(x, y)
			}
		}
	}
}

func BenchmarkLines_LineForWrite(b *testing.B) {
	var (
		color     = image.RGBA(10, 20, 30, 40)
		img       = image.New(acceleratedImageStub{width: 1920, height: 1080})
		selection = img.WholeImageSelection()
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lines := selection.Lines()
		for y := 0; y < lines.Length(); y++ {
			line := lines.LineForWrite(y)
			for x := 0; x < len(line); x++ {
				line[x] = color
			}
		}
	}
}

func BenchmarkLines_LineForRead(b *testing.B) {
	var (
		img       = image.New(acceleratedImageStub{width: 1920, height: 1080})
		selection = img.WholeImageSelection()
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lines := selection.Lines()
		for y := 0; y < lines.Length(); y++ {
			line := lines.LineForRead(y)
			for x := 0; x < len(line); x++ {
				_ = line[x]
			}
		}
	}
}

// Must be 0 allocs/op
func BenchmarkImage_Selection(b *testing.B) {
	var (
		img = image.New(acceleratedImageStub{width: 1, height: 1})
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img.Selection(1, 2).WithSize(3, 4)
	}
}

// Must be 0 allocs/op
func BenchmarkSelection_Selection(b *testing.B) {
	var (
		img       = image.New(acceleratedImageStub{width: 1, height: 1})
		selection = img.WholeImageSelection()
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selection.Selection(1, 2).WithSize(3, 4)
	}
}

// Must be 0 allocs/op
func BenchmarkSelection_Modify(b *testing.B) {
	var (
		img             = image.New(acceleratedImageStub{width: 1, height: 1})
		selection       = img.WholeImageSelection()
		command         = &acceleratedCommandStub{}
		sourceSelection = img.WholeImageSelection()
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selection.Modify(command, sourceSelection)
	}
}

type acceleratedCommandStub struct {
}

func (a acceleratedCommandStub) Run(output image.AcceleratedImageSelection, selections []image.AcceleratedImageSelection) {
}
