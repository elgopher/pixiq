package keyboard

// EventQueue is a data structure which can be used by libraries
// to collect keyboard Events in memory. EventQueue is an EventSource and can be
// directly consumed by Keyboard.
type EventQueue struct {
	circularBuffer []Event
	writeIndex     int
	readIndex      int
	available      int
}

// NewEventQueue creates EventQueue of given size. The minimum size of queue is 1.
// Size smaller than 1 is constrained to 1.
func NewEventQueue(size int) *EventQueue {
	if size < 1 {
		size = 1
	}
	return &EventQueue{circularBuffer: make([]Event, size)}
}

// Append add event to the queue. If there is not enough space the oldest event
// will be replaced.
func (q *EventQueue) Append(event Event) {
	if len(q.circularBuffer) == q.writeIndex {
		q.writeIndex = 0
		if q.readIndex == 0 {
			q.readIndex++
		}
		q.available--
	}
	q.circularBuffer[q.writeIndex] = event
	q.writeIndex++
	q.available++
}

// Poll retrieves and removes event from the queue. If there are no available
// events EmptyEvent and false is returned.
func (q *EventQueue) Poll() (Event, bool) {
	if q.available == 0 {
		return EmptyEvent, false
	}
	if len(q.circularBuffer) == q.readIndex {
		q.readIndex = 0
	}
	event := q.circularBuffer[q.readIndex]
	q.readIndex++
	q.available--
	return event, true
}
