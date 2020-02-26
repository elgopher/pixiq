package opengl

import (
	"bytes"
	"log"
	"runtime"
	"strconv"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// StartMainThreadLoop starts a loop assigned to main thread. It has to be
// executed from main goroutine or will panic. This function takes control over
// current goroutine by blocking it until runInDifferentGoroutine finishes.
// It provides Execute() method which can be used to execute given piece of code
// inside the main thread.
func StartMainThreadLoop(runInDifferentGoroutine func(*MainThreadLoop)) {
	if !isMainGoroutine() {
		panic("opengl.StartMainThreadLoop must be executed from main goroutine")
	}
	runtime.LockOSThread()
	commands := make(chan command, 4096)
	loop := &MainThreadLoop{
		commands: commands,
	}
	go func() {
		runInDifferentGoroutine(loop)
		close(commands)
	}()
	loop.run()
}

func isMainGoroutine() bool {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n == 1
}

// MainThreadLoop is a loop for executing jobs in main thread.
type MainThreadLoop struct {
	commands    chan command
	boundWindow *glfw.Window
}

func (g *MainThreadLoop) run() {
	defer logPanic()
	for {
		cmd, ok := <-g.commands
		if !ok {
			return
		}
		if cmd.window != nil {
			g.bind(cmd.window)
		}
		cmd.execute()
	}
}

func logPanic() {
	if p := recover(); p != nil {
		log.Panicln("panic in main thread loop: ", p)
	}
}

// Execute runs job blocking the current goroutine.
func (g *MainThreadLoop) Execute(job func()) {
	if job == nil {
		return
	}
	done := make(chan struct{})
	g.commands <- command{
		execute: func() {
			job()
			done <- struct{}{}
		},
	}
	<-done
}

func (g *MainThreadLoop) bind(window *glfw.Window) {
	if g.boundWindow != window {
		window.MakeContextCurrent()
		g.boundWindow = window
	}
}

// better use an array and create a function which will decode arguments
type command struct {
	window  *glfw.Window
	execute func()
}

func (g *MainThreadLoop) executeAsyncCommand(command command) {
	g.commands <- command
}

func (g *MainThreadLoop) executeCommand(cmd command) {
	done := make(chan struct{})
	g.commands <- command{
		window: cmd.window,
		execute: func() {
			cmd.execute()
			done <- struct{}{}
		},
	}
	<-done
}
