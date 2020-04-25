package image

import "fmt"

// Transparent is a special color where each component is zero (including alpha)
var Transparent = RGBA(0, 0, 0, 0)

// Color represents pixel color using 4 components: Red, Green, Black and Alpha.
// Red, Green and Blue components are premultiplied by alpha.
//
// Color is immutable struct. Changing the color means creating a new instance.
type Color struct {
	r, g, b, a byte
}

// R returns the red component.
func (c Color) R() byte {
	return c.r
}

// G returns the green component.
func (c Color) G() byte {
	return c.g
}

// B returns the blue component.
func (c Color) B() byte {
	return c.b
}

// A returns the alpha component.
func (c Color) A() byte {
	return c.a
}

// RGBA returns all color components as bytes in range 0 to 255.
func (c Color) RGBA() (byte, byte, byte, byte) {
	return c.r, c.g, c.b, c.a
}

// RGBAi returns all color components as integers in range 0 to 255.
func (c Color) RGBAi() (int, int, int, int) {
	return int(c.r), int(c.g), int(c.b), int(c.a)
}

// RGBAf returns all color components as floats in range 0.0 to 1.0.
func (c Color) RGBAf() (float32, float32, float32, float32) {
	return float32(c.r) / 255.0,
		float32(c.g) / 255.0,
		float32(c.b) / 255.0,
		float32(c.a) / 255.0
}

func (c Color) String() string {
	return fmt.Sprintf("RGBA(%d, %d, %d, %d)", c.r, c.g, c.b, c.a)
}

// RGBA creates Color using all four components: red, green, blue and alpha.
func RGBA(r, g, b, a byte) Color {
	return Color{
		r: r,
		g: g,
		b: b,
		a: a,
	}
}

// NRGBA creates Color using RGB components not premultiplied by alpha (aka straight
// alpha). Straight colors are being used by programs such as Aseprite.
func NRGBA(r, g, b, a byte) Color {
	return Color{
		r: mul(r, a),
		g: mul(g, a),
		b: mul(b, a),
		a: a,
	}
}

// mul is an optimized version of round(a * b / 255)
func mul(a, b byte) byte {
	t := int(a)*int(b) + 0x80
	return byte(((t >> 8) + t) >> 8)
}

// RGB creates Color using three components: red, green and blue.
// The color will be fully opaque (alpha=255)
func RGB(r, g, b byte) Color {
	return Color{
		r: r,
		g: g,
		b: b,
		a: 255,
	}
}

// RGBAi creates Color using components given as integer values.
// All values are clamped to [0-255] range.
func RGBAi(r, g, b, a int) Color {
	if r < 0 {
		r = 0
	}
	if r > 255 {
		r = 255
	}
	if g < 0 {
		g = 0
	}
	if g > 255 {
		g = 255
	}
	if b < 0 {
		b = 0
	}
	if b > 255 {
		b = 255
	}
	if a < 0 {
		a = 0
	}
	if a > 255 {
		a = 255
	}
	return Color{
		r: byte(r),
		g: byte(g),
		b: byte(b),
		a: byte(a),
	}
}
