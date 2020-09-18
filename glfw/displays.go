package glfw

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// Displays returns an object for getting information about
// currently connected displays.
//
// Displays instance requires MainThreadLoop because accessing information
// about displays must be done from the main thread.
func Displays(loop *MainThreadLoop) (*DisplaysAPI, error) {
	var err error
	loop.Execute(func() {
		err = glfw.Init()
	})
	if err != nil {
		return nil, fmt.Errorf("glfw.Init failed: %v", err)
	}
	return &DisplaysAPI{loop: loop}, nil

}

// DisplaysAPI provides information about currently connected displays.
type DisplaysAPI struct {
	loop *MainThreadLoop
}

// Primary returns the primary display. This is usually the display
// where elements like the Windows task bar or the OS X menu bar is located.
//
// Second return value is false when primary display does not exist (possibly
// no monitors are connected)
func (m *DisplaysAPI) Primary() (*Display, bool) {
	var monitor *glfw.Monitor
	m.loop.Execute(func() {
		monitor = glfw.GetPrimaryMonitor()
	})
	if monitor == nil {
		return nil, false
	}
	return &Display{
		monitor: monitor,
		loop:    m.loop,
	}, true
}

// All returns all connected displays
func (m *DisplaysAPI) All() []Display {
	var all []Display
	var glfwMonitors []*glfw.Monitor
	m.loop.Execute(func() {
		glfwMonitors = glfw.GetMonitors()
	})
	for _, monitor := range glfwMonitors {
		all = append(all,
			Display{
				monitor: monitor,
				loop:    m.loop,
			})
	}
	return all
}

// Display (aka monitor) provides information about display
type Display struct {
	monitor *glfw.Monitor
	loop    *MainThreadLoop
}

// Name returns a human-readable name of the display
func (m Display) Name() (name string) {
	m.loop.Execute(func() {
		name = m.monitor.GetName()
	})
	return
}

// Workarea returns the position, in pixels, of the upper-left
// corner of the work area of the specified display along with the work area
// size in pixels.
func (m Display) Workarea() (area Workarea) {
	m.loop.Execute(func() {
		x, y, width, height := m.monitor.GetWorkarea()
		area = Workarea{
			x:      x,
			y:      y,
			width:  width,
			height: height,
		}
	})
	return
}

// Workarea is the position, in pixels, of the upper-left
// corner of the work area of the specified display along with the work area
// size in pixels. The work area is defined as the area of the
// monitor not occluded by the operating system task bar where present. If no
// task bar exists then the work area is the monitor resolution in screen
// coordinates.
type Workarea struct {
	x, y, width, height int
}

// X coordinate of the upper-left corner of the work area in pixels
func (w Workarea) X() int {
	return w.x
}

// Y coordinate of the upper-left corner of the work area in pixels
func (w Workarea) Y() int {
	return w.y
}

// Width in pixels
func (w Workarea) Width() int {
	return w.width
}

// Height in pixels
func (w Workarea) Height() int {
	return w.height
}

// VideoMode returns the current video mode of the display
func (m Display) VideoMode() VideoMode {
	var mode *glfw.VidMode
	m.loop.Execute(func() {
		mode = m.monitor.GetVideoMode()
	})
	if mode == nil {
		panic("nil mode")
	}
	return VideoMode{
		width:       mode.Width,
		height:      mode.Height,
		refreshRate: mode.RefreshRate,
		monitor:     m.monitor,
	}
}

// VideoModes returns all video modes supported by the display.
// The returned array is sorted in ascending order by resolution area
// (the product of width and height).
func (m Display) VideoModes() []VideoMode {
	var modes []*glfw.VidMode
	m.loop.Execute(func() {
		modes = m.monitor.GetVideoModes()
	})
	var videoModes []VideoMode
	for _, mode := range modes {
		videoModes = append(videoModes, VideoMode{
			width:       mode.Width,
			height:      mode.Height,
			refreshRate: mode.RefreshRate,
			monitor:     m.monitor,
		})
	}
	return videoModes
}

// VideoMode contains information about display resolution in pixels
type VideoMode struct {
	monitor     *glfw.Monitor
	width       int
	height      int
	refreshRate int
}

// Width in pixels
func (m VideoMode) Width() int {
	return m.width
}

// Height in pixels
func (m VideoMode) Height() int {
	return m.height
}

// RefreshRate in hertz
func (m VideoMode) RefreshRate() int {
	return m.refreshRate
}

// PhysicalSize returns the size of the display area of the monitor.
func (m Display) PhysicalSize() (size PhysicalSize) {
	m.loop.Execute(func() {
		w, h := m.monitor.GetPhysicalSize()
		size = PhysicalSize{
			width:  w,
			height: h,
		}
	})
	return
}

// PhysicalSize returns the size, in millimetres, of the display area of the
// monitor.
//
// Note: Some operating systems do not provide accurate information, either
// because the monitor's EDID data is incorrect, or because the driver does not
// report it accurately.
type PhysicalSize struct {
	width  int
	height int
}

// Width in millimetres
func (s PhysicalSize) Width() int {
	return s.width
}

// Height in millimetres
func (s PhysicalSize) Height() int {
	return s.height
}
