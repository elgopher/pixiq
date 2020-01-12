package keyboard

// EventBuffer is a capped collection of accumulated events which can
// be used by libraries or in unit tests as a fake implementation of EventSource.
// The order of added events is preserved.
// EventBuffer is an EventSource and can be directly consumed by Keyboard.
type EventBuffer struct {
	circularBuffer []Event
	writeIndex     int
	readIndex      int
	readAfterWrite bool
}

// NewEventBuffer creates EventBuffer of given size. The minimum size of buffer is 1.
// Size smaller than 1 is constrained to 1.
func NewEventBuffer(size int) *EventBuffer {
	if size < 1 {
		size = 1
	}
	return &EventBuffer{circularBuffer: make([]Event, size)}
}

// Add adds event to the buffer. If there is not enough space the oldest event
// will be replaced.
func (q *EventBuffer) Add(event Event) {
	if len(q.circularBuffer) == q.writeIndex {
		q.writeIndex = 0
		q.readAfterWrite = true
	}
	if q.readAfterWrite && q.readIndex == q.writeIndex {
		q.readIndex++
	}
	q.circularBuffer[q.writeIndex] = event
	q.writeIndex++
}

// Poll retrieves and removes event from the buffer. If there are no available
// events EmptyEvent and false is returned.
func (q *EventBuffer) Poll() (Event, bool) {
	if q.writeIndex == q.readIndex && !q.readAfterWrite {
		return EmptyEvent, false
	}
	if len(q.circularBuffer) == q.readIndex {
		q.readIndex = 0
		q.readAfterWrite = false
	}
	event := q.circularBuffer[q.readIndex]
	q.readIndex++
	return event, true
}
