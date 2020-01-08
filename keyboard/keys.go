package keyboard

var (
	// Space Key " "
	Space = newKey(" ")
	// Apostrophe Key "'"
	Apostrophe = newKey("'")
	// Comma Key ","
	Comma = newKey(",")
	// Minus Key "-"
	Minus = newKey("-")
	// Period Key "."
	Period = newKey(".")
	// Slash Key "/"
	Slash = newKey("/")
	// Zero Key "0"
	Zero = newKey("0")
	// One Key "1"
	One = newKey("1")
	// Two Key "2"
	Two = newKey("2")
	// Three Key "3"
	Three = newKey("3")
	// Four Key "4"
	Four = newKey("4")
	// Five Key "5"
	Five = newKey("5")
	// Six Key "6"
	Six = newKey("6")
	// Seven Key "7"
	Seven = newKey("7")
	// Eight Key "8"
	Eight = newKey("8")
	// Nine Key "9"
	Nine = newKey("9")
	// Semicolon Key ";"
	Semicolon = newKey(";")
	// Equal Key "="
	Equal = newKey("=")
	// A key
	A = newKey("A")
	// B key
	B = newKey("B")
	// C Key
	C = newKey("C")
	// D Key
	D = newKey("D")
	// E Key
	E = newKey("E")
	// F Key
	F = newKey("F")
	// G Key
	G = newKey("G")
	// H Key
	H = newKey("H")
	// I Key
	I = newKey("I")
	// J Key
	J = newKey("J")
	// K Key
	K = newKey("K")
	// L Key
	L = newKey("L")
	// M Key
	M = newKey("M")
	// N Key
	N = newKey("N")
	// O Key
	O = newKey("O")
	// P Key
	P = newKey("P")
	// Q Key
	Q = newKey("Q")
	// R Key
	R = newKey("R")
	// S Key
	S = newKey("S")
	// T Key
	T = newKey("T")
	// U Key
	U = newKey("U")
	// V Key
	V = newKey("V")
	// W Key
	W = newKey("W")
	// X Key
	X = newKey("X")
	// Y Key
	Y = newKey("Y")
	// Z Key
	Z = newKey("Z")
	// LeftBracket Key "["
	LeftBracket = newKey("[")
	// Backslash Key "\"
	Backslash = newKey("\\")
	// RightBracket Key "]"
	RightBracket = newKey("]")
	// GraveAccent Key "`"
	GraveAccent = newKey("`")
	// World1 Key
	World1 = newKey("World 1")
	// World2 Key
	World2 = newKey("World 2")
	// Esc Key
	Esc = newKey("Esc")
	// Enter Key
	Enter = newKey("Enter")
	// Tab Key
	Tab = newKey("Tab")
	// Backspace Key
	Backspace = newKey("Backspace")
	// Insert Key
	Insert = newKey("Insert")
	// Delete Key
	Delete = newKey("Delete")
	// Right Key
	Right = newKey("Right")
	// Left Key
	Left = newKey("Left")
	// Down Key
	Down = newKey("Down")
	// Up Key
	Up = newKey("Up")
	// PageUp Key
	PageUp = newKey("PageUp")
	// PageDown Key
	PageDown = newKey("PageDown")
	// Home Key
	Home = newKey("Home")
	// End Key
	End = newKey("End")
	// CapsLock Key
	CapsLock = newKey("CapsLock")
	// ScrollLock Key
	ScrollLock = newKey("ScrollLock")
	// NumLock Key
	NumLock = newKey("NumLock")
	// PrintScreen Key
	PrintScreen = newKey("PrintScreen")
	// Pause Key
	Pause = newKey("Pause")
	// F1 Key
	F1 = newKey("F1")
	// F2 Key
	F2 = newKey("F2")
	// F3 Key
	F3 = newKey("F3")
	// F4 Key
	F4 = newKey("F4")
	// F5 Key
	F5 = newKey("F5")
	// F6 Key
	F6 = newKey("F6")
	// F7 Key
	F7 = newKey("F7")
	// F8 Key
	F8 = newKey("F8")
	// F9 Key
	F9 = newKey("F9")
	// F10 Key
	F10 = newKey("F10")
	// F11 Key
	F11 = newKey("F11")
	// F12 Key
	F12 = newKey("F12")
	// F13 Key
	F13 = newKey("F13")
	// F14 Key
	F14 = newKey("F14")
	// F15 Key
	F15 = newKey("F15")
	// F16 Key
	F16 = newKey("F16")
	// F17 Key
	F17 = newKey("F17")
	// F18 Key
	F18 = newKey("F18")
	// F19 Key
	F19 = newKey("F19")
	// F20 Key
	F20 = newKey("F20")
	// F21 Key
	F21 = newKey("F21")
	// F22 Key
	F22 = newKey("F22")
	// F23 Key
	F23 = newKey("F23")
	// F24 Key
	F24 = newKey("F24")
	// F25 Key
	F25 = newKey("F25")
	// Keypad0 Key
	Keypad0 = newKey("Keypad 0")
	// Keypad1 Key
	Keypad1 = newKey("Keypad 1")
	// Keypad2 Key
	Keypad2 = newKey("Keypad 2")
	// Keypad3 Key
	Keypad3 = newKey("Keypad 3")
	// Keypad4 Key
	Keypad4 = newKey("Keypad 4")
	// Keypad5 Key
	Keypad5 = newKey("Keypad 5")
	// Keypad6 Key
	Keypad6 = newKey("Keypad 6")
	// Keypad7 Key
	Keypad7 = newKey("Keypad 7")
	// Keypad8 Key
	Keypad8 = newKey("Keypad 8")
	// Keypad9 Key
	Keypad9 = newKey("Keypad 9")
	// KeypadDecimal Key
	KeypadDecimal = newKey("Keypad .")
	// KeypadDivide Key
	KeypadDivide = newKey("Keypad /")
	// KeypadMultiply Key
	KeypadMultiply = newKey("Keypad *")
	// KeypadSubtract Key
	KeypadSubtract = newKey("Keypad -")
	// KeypadAdd Key
	KeypadAdd = newKey("Keypad +")
	// KeypadEnter Key
	KeypadEnter = newKey("Keypad Enter")
	// KeypadEqual Key
	KeypadEqual = newKey("Keypad =")
	// LeftShift Key
	LeftShift = newKey("Left Shift")
	// LeftControl Key
	LeftControl = newKey("Left Control")
	// LeftAlt Key
	LeftAlt = newKey("Left Alt")
	// LeftSuper Key
	LeftSuper = newKey("Left Super")
	// RightShift Key
	RightShift = newKey("Right Shift")
	// RightControl Key
	RightControl = newKey("Right Control")
	// RightAlt Key
	RightAlt = newKey("Right Alt")
	// RightSuper Key
	RightSuper = newKey("Right Super")
	// Menu Key
	Menu = newKey("Menu")

	allKeys = []Key{
		Space, Apostrophe, Comma, Minus, Period, Slash,
		Zero, One, Two, Three, Four, Five, Six, Seven, Eight, Nine,
		Semicolon, Equal,
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z,
		LeftBracket, Backslash, RightBracket, GraveAccent,
		World1, World2,
		Esc, Enter, Tab, Backspace, Insert, Delete,
		Right, Left, Down, Up, PageUp, PageDown,
		Home, End,
		CapsLock, ScrollLock, NumLock,
		Pause,
		F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12,
		F13, F14, F15, F16, F17, F18, F19, F20, F21, F22, F23, F24, F25,
		Keypad0, Keypad1, Keypad2, Keypad3, Keypad4, Keypad5, Keypad6, Keypad7, Keypad8, Keypad9,
		KeypadDecimal, KeypadDivide, KeypadMultiply, KeypadMultiply,
		KeypadSubtract, KeypadAdd, KeypadEnter, KeypadEqual,
		LeftShift, LeftControl, LeftAlt, LeftSuper,
		RightShift, RightControl, RightAlt, RightSuper,
		Menu,
	}
)
