package keyboard_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/keyboard"
)

func BenchmarkKeyboard_Update(b *testing.B) {
	event := keyboard.NewPressedEvent(keyboard.A)
	source := &fixedEventsSource{events: []keyboard.Event{event}}
	keys := keyboard.New(source)
	for i := 0; i < b.N; i++ {
		keys.Update() // should be 0 allocs/op
	}
}

type fixedEventsSource struct {
	events []keyboard.Event
}

func (n *fixedEventsSource) Poll(output []keyboard.Event) int {
	for i, event := range n.events {
		output[i] = event
	}
	return len(n.events)
}
