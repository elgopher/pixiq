// Package colornames provides named colors as defined in the SVG 1.1 spec.
// The package is inspired by golang.org/x/image/colornames
//
// See http://www.w3.org/TR/SVG/types.html#ColorKeywords
// See https://github.com/golang/image/tree/master/colornames
package colornames

import (
	"github.com/elgopher/pixiq/image"
)

var (
	// Aliceblue is pixiq.RGB(240, 248, 255)
	Aliceblue = image.RGB(240, 248, 255)
	// Antiquewhite is pixiq.RGB(250, 235, 215)
	Antiquewhite = image.RGB(250, 235, 215)
	// Aqua is pixiq.RGB(0, 255, 255)
	Aqua = image.RGB(0, 255, 255)
	// Aquamarine is pixiq.RGB(127, 255, 212)
	Aquamarine = image.RGB(127, 255, 212)
	// Azure is pixiq.RGB(240, 255, 255)
	Azure = image.RGB(240, 255, 255)
	// Beige is pixiq.RGB(245, 245, 220)
	Beige = image.RGB(245, 245, 220)
	// Bisque is pixiq.RGB(255, 228, 196)
	Bisque = image.RGB(255, 228, 196)
	// Black is pixiq.RGB(0, 0, 0)
	Black = image.RGB(0, 0, 0)
	// Blanchedalmond is pixiq.RGB(255, 235, 205)
	Blanchedalmond = image.RGB(255, 235, 205)
	// Blue is pixiq.RGB(0, 0, 255)
	Blue = image.RGB(0, 0, 255)
	// Blueviolet is pixiq.RGB(138, 43, 226)
	Blueviolet = image.RGB(138, 43, 226)
	// Brown is pixiq.RGB(165, 42, 42)
	Brown = image.RGB(165, 42, 42)
	// Burlywood is pixiq.RGB(222, 184, 135)
	Burlywood = image.RGB(222, 184, 135)
	// Cadetblue is pixiq.RGB(95, 158, 160)
	Cadetblue = image.RGB(95, 158, 160)
	// Chartreuse is pixiq.RGB(127, 255, 0)
	Chartreuse = image.RGB(127, 255, 0)
	// Chocolate is pixiq.RGB(210, 105, 30)
	Chocolate = image.RGB(210, 105, 30)
	// Coral is pixiq.RGB(255, 127, 80)
	Coral = image.RGB(255, 127, 80)
	// Cornflowerblue is pixiq.RGB(100, 149, 237)
	Cornflowerblue = image.RGB(100, 149, 237)
	// Cornsilk is pixiq.RGB(255, 248, 220)
	Cornsilk = image.RGB(255, 248, 220)
	// Crimson is pixiq.RGB(220, 20, 60)
	Crimson = image.RGB(220, 20, 60)
	// Cyan is pixiq.RGB(0, 255, 255)
	Cyan = image.RGB(0, 255, 255)
	// Darkblue is pixiq.RGB(0, 0, 139)
	Darkblue = image.RGB(0, 0, 139)
	// Darkcyan is pixiq.RGB(0, 139, 139)
	Darkcyan = image.RGB(0, 139, 139)
	// Darkgoldenrod is pixiq.RGB(184, 134, 11)
	Darkgoldenrod = image.RGB(184, 134, 11)
	// Darkgray is pixiq.RGB(169, 169, 169)
	Darkgray = image.RGB(169, 169, 169)
	// Darkgreen is pixiq.RGB(0, 100, 0)
	Darkgreen = image.RGB(0, 100, 0)
	// Darkgrey is pixiq.RGB(169, 169, 169)
	Darkgrey = image.RGB(169, 169, 169)
	// Darkkhaki is pixiq.RGB(189, 183, 107)
	Darkkhaki = image.RGB(189, 183, 107)
	// Darkmagenta is pixiq.RGB(139, 0, 139)
	Darkmagenta = image.RGB(139, 0, 139)
	// Darkolivegreen is pixiq.RGB(85, 107, 47)
	Darkolivegreen = image.RGB(85, 107, 47)
	// Darkorange is pixiq.RGB(255, 140, 0)
	Darkorange = image.RGB(255, 140, 0)
	// Darkorchid is pixiq.RGB(153, 50, 204)
	Darkorchid = image.RGB(153, 50, 204)
	// Darkred is pixiq.RGB(139, 0, 0)
	Darkred = image.RGB(139, 0, 0)
	// Darksalmon is pixiq.RGB(233, 150, 122)
	Darksalmon = image.RGB(233, 150, 122)
	// Darkseagreen is pixiq.RGB(143, 188, 143)
	Darkseagreen = image.RGB(143, 188, 143)
	// Darkslateblue is pixiq.RGB(72, 61, 139)
	Darkslateblue = image.RGB(72, 61, 139)
	// Darkslategray is pixiq.RGB(47, 79, 79)
	Darkslategray = image.RGB(47, 79, 79)
	// Darkslategrey is pixiq.RGB(47, 79, 79)
	Darkslategrey = image.RGB(47, 79, 79)
	// Darkturquoise is pixiq.RGB(0, 206, 209)
	Darkturquoise = image.RGB(0, 206, 209)
	// Darkviolet is pixiq.RGB(148, 0, 211)
	Darkviolet = image.RGB(148, 0, 211)
	// Deeppink is pixiq.RGB(255, 20, 147)
	Deeppink = image.RGB(255, 20, 147)
	// Deepskyblue is pixiq.RGB(0, 191, 255)
	Deepskyblue = image.RGB(0, 191, 255)
	// Dimgray is pixiq.RGB(105, 105, 105)
	Dimgray = image.RGB(105, 105, 105)
	// Dimgrey is pixiq.RGB(105, 105, 105)
	Dimgrey = image.RGB(105, 105, 105)
	// Dodgerblue is pixiq.RGB(30, 144, 255)
	Dodgerblue = image.RGB(30, 144, 255)
	// Firebrick is pixiq.RGB(178, 34, 34)
	Firebrick = image.RGB(178, 34, 34)
	// Floralwhite is pixiq.RGB(255, 250, 240)
	Floralwhite = image.RGB(255, 250, 240)
	// Forestgreen is pixiq.RGB(34, 139, 34)
	Forestgreen = image.RGB(34, 139, 34)
	// Fuchsia is pixiq.RGB(255, 0, 255)
	Fuchsia = image.RGB(255, 0, 255)
	// Gainsboro is pixiq.RGB(220, 220, 220)
	Gainsboro = image.RGB(220, 220, 220)
	// Ghostwhite is pixiq.RGB(248, 248, 255)
	Ghostwhite = image.RGB(248, 248, 255)
	// Gold is pixiq.RGB(255, 215, 0)
	Gold = image.RGB(255, 215, 0)
	// Goldenrod is pixiq.RGB(218, 165, 32)
	Goldenrod = image.RGB(218, 165, 32)
	// Gray is pixiq.RGB(128, 128, 128)
	Gray = image.RGB(128, 128, 128)
	// Green is pixiq.RGB(0, 128, 0)
	Green = image.RGB(0, 128, 0)
	// Greenyellow is pixiq.RGB(173, 255, 47)
	Greenyellow = image.RGB(173, 255, 47)
	// Grey is pixiq.RGB(128, 128, 128)
	Grey = image.RGB(128, 128, 128)
	// Honeydew is pixiq.RGB(240, 255, 240)
	Honeydew = image.RGB(240, 255, 240)
	// Hotpink is pixiq.RGB(255, 105, 180)
	Hotpink = image.RGB(255, 105, 180)
	// Indianred is pixiq.RGB(205, 92, 92)
	Indianred = image.RGB(205, 92, 92)
	// Indigo is pixiq.RGB(75, 0, 130)
	Indigo = image.RGB(75, 0, 130)
	// Ivory is pixiq.RGB(255, 255, 240)
	Ivory = image.RGB(255, 255, 240)
	// Khaki is pixiq.RGB(240, 230, 140)
	Khaki = image.RGB(240, 230, 140)
	// Lavender is pixiq.RGB(230, 230, 250)
	Lavender = image.RGB(230, 230, 250)
	// Lavenderblush is pixiq.RGB(255, 240, 245)
	Lavenderblush = image.RGB(255, 240, 245)
	// Lawngreen is pixiq.RGB(124, 252, 0)
	Lawngreen = image.RGB(124, 252, 0)
	// Lemonchiffon is pixiq.RGB(255, 250, 205)
	Lemonchiffon = image.RGB(255, 250, 205)
	// Lightblue is pixiq.RGB(173, 216, 230)
	Lightblue = image.RGB(173, 216, 230)
	// Lightcoral is pixiq.RGB(240, 128, 128)
	Lightcoral = image.RGB(240, 128, 128)
	// Lightcyan is pixiq.RGB(224, 255, 255)
	Lightcyan = image.RGB(224, 255, 255)
	// Lightgoldenrodyellow is pixiq.RGB(250, 250, 210)
	Lightgoldenrodyellow = image.RGB(250, 250, 210)
	// Lightgray is pixiq.RGB(211, 211, 211)
	Lightgray = image.RGB(211, 211, 211)
	// Lightgreen is pixiq.RGB(144, 238, 144)
	Lightgreen = image.RGB(144, 238, 144)
	// Lightgrey is pixiq.RGB(211, 211, 211)
	Lightgrey = image.RGB(211, 211, 211)
	// Lightpink is pixiq.RGB(255, 182, 193)
	Lightpink = image.RGB(255, 182, 193)
	// Lightsalmon is pixiq.RGB(255, 160, 122)
	Lightsalmon = image.RGB(255, 160, 122)
	// Lightseagreen is pixiq.RGB(32, 178, 170)
	Lightseagreen = image.RGB(32, 178, 170)
	// Lightskyblue is pixiq.RGB(135, 206, 250)
	Lightskyblue = image.RGB(135, 206, 250)
	// Lightslategray is pixiq.RGB(119, 136, 153)
	Lightslategray = image.RGB(119, 136, 153)
	// Lightslategrey is pixiq.RGB(119, 136, 153)
	Lightslategrey = image.RGB(119, 136, 153)
	// Lightsteelblue is pixiq.RGB(176, 196, 222)
	Lightsteelblue = image.RGB(176, 196, 222)
	// Lightyellow is pixiq.RGB(255, 255, 224)
	Lightyellow = image.RGB(255, 255, 224)
	// Lime is pixiq.RGB(0, 255, 0)
	Lime = image.RGB(0, 255, 0)
	// Limegreen is pixiq.RGB(50, 205, 50)
	Limegreen = image.RGB(50, 205, 50)
	// Linen is pixiq.RGB(250, 240, 230)
	Linen = image.RGB(250, 240, 230)
	// Magenta is pixiq.RGB(255, 0, 255)
	Magenta = image.RGB(255, 0, 255)
	// Maroon is pixiq.RGB(128, 0, 0)
	Maroon = image.RGB(128, 0, 0)
	// Mediumaquamarine is pixiq.RGB(102, 205, 170)
	Mediumaquamarine = image.RGB(102, 205, 170)
	// Mediumblue is pixiq.RGB(0, 0, 205)
	Mediumblue = image.RGB(0, 0, 205)
	// Mediumorchid is pixiq.RGB(186, 85, 211)
	Mediumorchid = image.RGB(186, 85, 211)
	// Mediumpurple is pixiq.RGB(147, 112, 219)
	Mediumpurple = image.RGB(147, 112, 219)
	// Mediumseagreen is pixiq.RGB(60, 179, 113)
	Mediumseagreen = image.RGB(60, 179, 113)
	// Mediumslateblue is pixiq.RGB(123, 104, 238)
	Mediumslateblue = image.RGB(123, 104, 238)
	// Mediumspringgreen is pixiq.RGB(0, 250, 154)
	Mediumspringgreen = image.RGB(0, 250, 154)
	// Mediumturquoise is pixiq.RGB(72, 209, 204)
	Mediumturquoise = image.RGB(72, 209, 204)
	// Mediumvioletred is pixiq.RGB(199, 21, 133)
	Mediumvioletred = image.RGB(199, 21, 133)
	// Midnightblue is pixiq.RGB(25, 25, 112)
	Midnightblue = image.RGB(25, 25, 112)
	// Mintcream is pixiq.RGB(245, 255, 250)
	Mintcream = image.RGB(245, 255, 250)
	// Mistyrose is pixiq.RGB(255, 228, 225)
	Mistyrose = image.RGB(255, 228, 225)
	// Moccasin is pixiq.RGB(255, 228, 181)
	Moccasin = image.RGB(255, 228, 181)
	// Navajowhite is pixiq.RGB(255, 222, 173)
	Navajowhite = image.RGB(255, 222, 173)
	// Navy is pixiq.RGB(0, 0, 128)
	Navy = image.RGB(0, 0, 128)
	// Oldlace is pixiq.RGB(253, 245, 230)
	Oldlace = image.RGB(253, 245, 230)
	// Olive is pixiq.RGB(128, 128, 0)
	Olive = image.RGB(128, 128, 0)
	// Olivedrab is pixiq.RGB(107, 142, 35)
	Olivedrab = image.RGB(107, 142, 35)
	// Orange is pixiq.RGB(255, 165, 0)
	Orange = image.RGB(255, 165, 0)
	// Orangered is pixiq.RGB(255, 69, 0)
	Orangered = image.RGB(255, 69, 0)
	// Orchid is pixiq.RGB(218, 112, 214)
	Orchid = image.RGB(218, 112, 214)
	// Palegoldenrod is pixiq.RGB(238, 232, 170)
	Palegoldenrod = image.RGB(238, 232, 170)
	// Palegreen is pixiq.RGB(152, 251, 152)
	Palegreen = image.RGB(152, 251, 152)
	// Paleturquoise is pixiq.RGB(175, 238, 238)
	Paleturquoise = image.RGB(175, 238, 238)
	// Palevioletred is pixiq.RGB(219, 112, 147)
	Palevioletred = image.RGB(219, 112, 147)
	// Papayawhip is pixiq.RGB(255, 239, 213)
	Papayawhip = image.RGB(255, 239, 213)
	// Peachpuff is pixiq.RGB(255, 218, 185)
	Peachpuff = image.RGB(255, 218, 185)
	// Peru is pixiq.RGB(205, 133, 63)
	Peru = image.RGB(205, 133, 63)
	// Pink is pixiq.RGB(255, 192, 203)
	Pink = image.RGB(255, 192, 203)
	// Plum is pixiq.RGB(221, 160, 221)
	Plum = image.RGB(221, 160, 221)
	// Powderblue is pixiq.RGB(176, 224, 230)
	Powderblue = image.RGB(176, 224, 230)
	// Purple is pixiq.RGB(128, 0, 128)
	Purple = image.RGB(128, 0, 128)
	// Red is pixiq.RGB(255, 0, 0)
	Red = image.RGB(255, 0, 0)
	// Rosybrown is pixiq.RGB(188, 143, 143)
	Rosybrown = image.RGB(188, 143, 143)
	// Royalblue is pixiq.RGB(65, 105, 225)
	Royalblue = image.RGB(65, 105, 225)
	// Saddlebrown is pixiq.RGB(139, 69, 19)
	Saddlebrown = image.RGB(139, 69, 19)
	// Salmon is pixiq.RGB(250, 128, 114)
	Salmon = image.RGB(250, 128, 114)
	// Sandybrown is pixiq.RGB(244, 164, 96)
	Sandybrown = image.RGB(244, 164, 96)
	// Seagreen is pixiq.RGB(46, 139, 87)
	Seagreen = image.RGB(46, 139, 87)
	// Seashell is pixiq.RGB(255, 245, 238)
	Seashell = image.RGB(255, 245, 238)
	// Sienna is pixiq.RGB(160, 82, 45)
	Sienna = image.RGB(160, 82, 45)
	// Silver is pixiq.RGB(192, 192, 192)
	Silver = image.RGB(192, 192, 192)
	// Skyblue is pixiq.RGB(135, 206, 235)
	Skyblue = image.RGB(135, 206, 235)
	// Slateblue is pixiq.RGB(106, 90, 205)
	Slateblue = image.RGB(106, 90, 205)
	// Slategray is pixiq.RGB(112, 128, 144)
	Slategray = image.RGB(112, 128, 144)
	// Slategrey is pixiq.RGB(112, 128, 144)
	Slategrey = image.RGB(112, 128, 144)
	// Snow is pixiq.RGB(255, 250, 250)
	Snow = image.RGB(255, 250, 250)
	// Springgreen is pixiq.RGB(0, 255, 127)
	Springgreen = image.RGB(0, 255, 127)
	// Steelblue is pixiq.RGB(70, 130, 180)
	Steelblue = image.RGB(70, 130, 180)
	// Tan is pixiq.RGB(210, 180, 140)
	Tan = image.RGB(210, 180, 140)
	// Teal is pixiq.RGB(0, 128, 128)
	Teal = image.RGB(0, 128, 128)
	// Thistle is pixiq.RGB(216, 191, 216)
	Thistle = image.RGB(216, 191, 216)
	// Tomato is pixiq.RGB(255, 99, 71)
	Tomato = image.RGB(255, 99, 71)
	// Turquoise is pixiq.RGB(64, 224, 208)
	Turquoise = image.RGB(64, 224, 208)
	// Violet is pixiq.RGB(238, 130, 238)
	Violet = image.RGB(238, 130, 238)
	// Wheat is pixiq.RGB(245, 222, 179)
	Wheat = image.RGB(245, 222, 179)
	// White is pixiq.RGB(255, 255, 255)
	White = image.RGB(255, 255, 255)
	// Whitesmoke is pixiq.RGB(245, 245, 245)
	Whitesmoke = image.RGB(245, 245, 245)
	// Yellow is pixiq.RGB(255, 255, 0)
	Yellow = image.RGB(255, 255, 0)
	// Yellowgreen is pixiq.RGB(154, 205, 50)
	Yellowgreen = image.RGB(154, 205, 50)
)
