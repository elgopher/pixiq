package keyboard_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/keyboard"
)

func BenchmarkKeyboard_Update(b *testing.B) {
	var (
		event  = keyboard.NewPressedEvent(keyboard.A)
		source = &cyclicEventsSoure{event: event}
		keys   = keyboard.New(source)
	)
	for i := 0; i < b.N; i++ {
		keys.Update() // should be 0 allocs/op
	}
}

type cyclicEventsSoure struct {
	hasEvent bool
	event    keyboard.Event
}

func (f *cyclicEventsSoure) Poll() (keyboard.Event, bool) {
	f.hasEvent = !f.hasEvent
	if f.hasEvent {
		return f.event, true
	}
	return keyboard.EmptyEvent, false
}
