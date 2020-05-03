package mouse_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/mouse"
)

func BenchmarkMouse_Update(b *testing.B) {
	var (
		event      = mouse.NewPressedEvent(mouse.Left)
		source     = &cyclicEventsSource{event: event}
		mouseState = mouse.New(source)
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mouseState.Update() // should be 0 alloc/op
	}
}

func BenchmarkMouse_PressedButtons(b *testing.B) {
	var (
		event      = mouse.NewPressedEvent(mouse.Left)
		source     = &cyclicEventsSource{event: event}
		mouseState = mouse.New(source)
	)
	mouseState.Update()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mouseState.PressedButtons() // should be at most 1 alloc/op
	}
}

type cyclicEventsSource struct {
	hasEvent bool
	event    mouse.Event
}

func (f *cyclicEventsSource) PollMouseEvent() (mouse.Event, bool) {
	f.hasEvent = !f.hasEvent
	if f.hasEvent {
		return f.event, true
	}
	return mouse.EmptyEvent, false
}
