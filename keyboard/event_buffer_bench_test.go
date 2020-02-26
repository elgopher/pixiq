package keyboard_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/keyboard"
)

// Should be 0 allocs/op
func BenchmarkKeyboardEvents(b *testing.B) {
	const size = 8
	events := keyboard.NewEventBuffer(size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < size*2; j++ {
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
