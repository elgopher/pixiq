package mouse_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/mouse"
)

// Should be 0 allocs/op
func BenchmarkMouseEvents(b *testing.B) {
	const size = 8
	events := mouse.NewEventBuffer(size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < size*2; j++ {
			events.Add(mouse.NewReleasedEvent(mouse.Left))
			events.Add(mouse.NewPressedEvent(mouse.Left))
			events.Add(mouse.NewScrolledEvent(1, 2))
			events.Add(mouse.NewMovedEvent(1, 2, 1, 2, true))
		}
		for {
			_, ok := events.Poll()
			if !ok {
				break
			}
		}
	}
}
