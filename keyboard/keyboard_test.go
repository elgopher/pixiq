package keyboard_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/keyboard"
)

func TestNew(t *testing.T) {
	t.Run("should panic when source is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			keyboard.New(nil)
		})
	})
	t.Run("should create a keyboard instance", func(t *testing.T) {
		source := &fakeEventsSource{}
		// when
		keys := keyboard.New(source)
		// then
		assert.NotNil(t, keys)
	})
}

func TestNewKey(t *testing.T) {
	t.Run("should panic for invalid tokens", func(t *testing.T) {
		tests := []keyboard.Token{1, 64, 67}
		for _, token := range tests {
			testName := fmt.Sprintf("for %v", token)
			t.Run(testName, func(t *testing.T) {
				assert.Panics(t, func() {
					keyboard.NewKey(token, 0)
				})
			})
		}
	})
	t.Run("should create new key", func(t *testing.T) {
		tests := map[string]struct {
			token    keyboard.Token
			scanCode int
		}{
			"A 0": {
				token:    keyboard.A,
				scanCode: 0,
			},
			"B 1": {
				token:    keyboard.B,
				scanCode: 1,
			},
			"Unknown 1": {
				token:    keyboard.Unknown,
				scanCode: 2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				key := keyboard.NewKey(test.token, test.scanCode)
				assert.Equal(t, test.token, key.Token())
				assert.Equal(t, test.scanCode, key.ScanCode())
			})
		}
	})
}

func TestNewPressedEvent(t *testing.T) {
	t.Run("should create PressedEvent", func(t *testing.T) {
		key := keyboard.NewKey(keyboard.A, 1)
		// when
		keyboard.NewPressedEvent(key)
	})
}

func TestKeyboard_Pressed(t *testing.T) {
	t.Run("before Update was called, Pressed returns false for all tokens", func(t *testing.T) {
		tests := []keyboard.Token{keyboard.A, keyboard.B}
		for _, token := range tests {
			testName := fmt.Sprintf("for token: %v", token)
			t.Run(testName, func(t *testing.T) {
				source := newFakeEventSourceWithPressedEvents(keyboard.A)
				keys := keyboard.New(source)
				// when
				pressed := keys.Pressed(token)
				// then
				assert.False(t, pressed)
			})
		}
	})
	t.Run("after Update was called", func(t *testing.T) {
		tests := map[string]struct {
			source             keyboard.EventsSource
			expectedPressed    []keyboard.Token
			expectedNotPressed []keyboard.Token
		}{
			"one event": {
				source:             newFakeEventSourceWithPressedEvents(keyboard.A),
				expectedPressed:    []keyboard.Token{keyboard.A},
				expectedNotPressed: []keyboard.Token{keyboard.B},
			},
			"two events for tokens B and A": {
				source:          newFakeEventSourceWithPressedEvents(keyboard.B, keyboard.A),
				expectedPressed: []keyboard.Token{keyboard.A, keyboard.B},
			},
			"two events for tokens A and B": {
				source:          newFakeEventSourceWithPressedEvents(keyboard.A, keyboard.B),
				expectedPressed: []keyboard.Token{keyboard.A, keyboard.B},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				keys := keyboard.New(test.source)
				// when
				keys.Update()
				// then
				for _, expectedPressedToken := range test.expectedPressed {
					assert.True(t, keys.Pressed(expectedPressedToken))
				}
				for _, expectedNotPressedToken := range test.expectedNotPressed {
					assert.False(t, keys.Pressed(expectedNotPressedToken))
				}
			})
		}
	})
}

func newFakeEventSourceWithPressedEvents(token ...keyboard.Token) *fakeEventsSource {
	source := &fakeEventsSource{}
	source.events = []keyboard.Event{}
	for _, token := range token {
		key := keyboard.NewKey(token, 1)
		event := keyboard.NewPressedEvent(key)
		source.events = append(source.events, event)
	}
	return source
}

type fakeEventsSource struct {
	events []keyboard.Event
}

func (f *fakeEventsSource) Poll(output []keyboard.Event) int {
	to := len(f.events)
	if len(output) < to {
		to = len(output)
	}
	slice := f.events[:to]
	for i, event := range slice {
		output[i] = event
	}
	f.events = []keyboard.Event{}
	return to
}
