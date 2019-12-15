package pixiq

// Color represents pixel color using 4 components: Red, Green, Black and Alpha.
// Red, Green and Blue components are not premultiplied by alpha (aka straight alpha), that is RGB and alpha are
// independent. You can change one without affecting the other.
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

// RGBA creates Color using all four components: red, green, blue and alpha.
func RGBA(r, g, b, a byte) Color {
	return Color{
		r: r,
		g: g,
		b: b,
		a: a,
	}
}

// RGBAi creates Color using components given as integer values. All values are clamped to [0-255] range.
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
