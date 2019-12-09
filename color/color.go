package color

// Color represents pixel color described by 4 components: red (r), green (g), blue (b) and alpha (a).
// Non-premultiplied alpha is used.
type Color struct {
	r, g, b, a uint8
}

// RGBA creates a color from 4 component. It is expected that all 4 components are in 0-255 range.
func RGBA(r, g, b, a int) Color {
	return Color{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// R returns Red component disregarding the opacity (non-premultiplied alpha aka straight alpha representation)
func (c Color) R() int {
	return int(c.r)
}

// G returns Green component disregarding the opacity (non-premultiplied alpha aka straight alpha representation)
func (c Color) G() int {
	return int(c.g)
}

// B returns Blue component disregarding the opacity (non-premultiplied alpha aka straight alpha representation)
func (c Color) B() int {
	return int(c.b)
}

// A returns Alpha component (aka opacity)
func (c Color) A() int {
	return int(c.a)
}
