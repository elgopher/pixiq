package opengl

import (
	"bytes"
	"log"
	"runtime"
	"strconv"
	"sync"

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
	jobs := make(chan func(), 4096)
	synchronousJob := &synchronousJob{
		done: make(chan struct{}),
	}
	loop := &MainThreadLoop{
		jobs:           jobs,
		synchronousJob: synchronousJob,
		runSynchronousJob: func() {
			synchronousJob.job()
			synchronousJob.done <- struct{}{}
		},
	}
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

// MainThreadLoop is a loop for executing jobs in main thread.
type MainThreadLoop struct {
	synchronousJob    *synchronousJob
	runSynchronousJob func()
	jobs              chan func()
	boundWindow       *glfw.Window
}

type synchronousJob struct {
	sync.Mutex
	job  func()
	done chan struct{}
}

func (g *MainThreadLoop) run() {
	defer logPanic()
	for {
		job, ok := <-g.jobs
		if !ok {
			return
		}
		job()
	}
}

func logPanic() {
	if p := recover(); p != nil {
		log.Panicln("panic in main thread loop: ", p)
	}
}

// Execute runs job blocking the current goroutine.
func (g *MainThreadLoop) Execute(job func()) {
	g.synchronousJob.Lock()
	defer g.synchronousJob.Unlock()
	g.synchronousJob.job = job
	g.jobs <- g.runSynchronousJob
	<-g.synchronousJob.done
}

// ExecuteAsync runs job asynchronously.
func (g *MainThreadLoop) ExecuteAsync(job func()) {
	g.jobs <- job
}

func (g *MainThreadLoop) bind(window *glfw.Window) {
	if g.boundWindow != window {
		window.MakeContextCurrent()
		g.boundWindow = window
	}
}
