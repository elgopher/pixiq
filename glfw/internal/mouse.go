package internal

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type MouseEvents struct {
}

func (e *MouseEvents) OnMouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {

}

// This is not fired when cursor is outside the window (Linux)
func (e *MouseEvents) OnCursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	//fmt.Println(xpos, ypos)
}

func (e *MouseEvents) OnCursorEnterCallback(w *glfw.Window, entered bool) {
	//fmt.Println(entered)
}

func (e *MouseEvents) OnScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	fmt.Println(xoff, yoff)
}
