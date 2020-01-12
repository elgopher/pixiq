package keyboard_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/keyboard"
)

// Should be 0 allocs/op
func BenchmarkKeyboardEvents(b *testing.B) {
	b.StopTimer()
	const size = 8
	events := keyboard.NewEventQueue(size)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < size; i++ {
			events.Append(keyboard.NewPressedEvent(keyboard.A))
		}
		for {
			_, ok := events.Poll()
			if !ok {
				break
			}
		}
	}
}
