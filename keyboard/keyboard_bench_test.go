package keyboard_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/keyboard"
)

func BenchmarkKeyboard_Update(b *testing.B) {
	var (
		event  = keyboard.NewPressedEvent(keyboard.A)
		source = &cyclicEventsSource{event: event}
		keys   = keyboard.New(source)
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		keys.Update() // should be 0 alloc/op
	}
}

func BenchmarkKeyboard_PressedKeys(b *testing.B) {
	var (
		event  = keyboard.NewPressedEvent(keyboard.A)
		source = &cyclicEventsSource{event: event}
		keys   = keyboard.New(source)
	)
	keys.Update()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		keys.PressedKeys() // should be at most 1 alloc/op
	}
}

type cyclicEventsSource struct {
	hasEvent bool
	event    keyboard.Event
}

func (f *cyclicEventsSource) PollKeyboardEvent() (keyboard.Event, bool) {
	f.hasEvent = !f.hasEvent
	if f.hasEvent {
		return f.event, true
	}
	return keyboard.EmptyEvent, false
}
