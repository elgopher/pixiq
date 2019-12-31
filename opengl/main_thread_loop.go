package opengl

import (
	"bytes"
	"runtime"
	"strconv"
)

// StartMainThreadLoop starts a main loop assigned to main thread. It has to be executed from main goroutine or will panic.
// This function takes control over current goroutine by blocking it until runInDifferentGoroutine finishes.
// It provides Execute() method which can be used to execute given piece of code inside the main thread.
func StartMainThreadLoop(runInDifferentGoroutine func(*MainThreadLoop)) {
	if !isMainGoroutine() {
		panic("opengl.StartMainThreadLoop must be executed from main goroutine")
	}
	runtime.LockOSThread()
	jobs := make(chan func())
	loop := &MainThreadLoop{jobs: jobs}
	go func() {
		runInDifferentGoroutine(loop)
		close(jobs)
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

// MainThreadLoop is a loop for executing jobs in main thread
type MainThreadLoop struct {
	jobs chan func()
}

func (g *MainThreadLoop) run() {
	for {
		job, ok := <-g.jobs
		if !ok {
			return
		}
		job()
	}
}

// Execute runs job blocking the current goroutine
func (g *MainThreadLoop) Execute(job func()) {
	done := make(chan struct{})
	g.jobs <- func() {
		job()
		done <- struct{}{}
	}
	<-done
}
