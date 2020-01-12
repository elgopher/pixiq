package keyboard_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/keyboard"
)

// Should be 0 allocs/op
func BenchmarkKeyboardEvents(b *testing.B) {
	b.StopTimer()
	const size = 8
	events := keyboard.NewEventBuffer(size)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < size*2; i++ {
			events.Add(keyboard.NewPressedEvent(keyboard.A))
		}
		for {
			_, ok := events.Poll()
			if !ok {
				break
			}
		}
	}
}
