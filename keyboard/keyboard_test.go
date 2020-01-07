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
		source := &fakeEventSource{}
		// when
		keys := keyboard.New(source)
		// then
		assert.NotNil(t, keys)
	})
}

func TestNewKey(t *testing.T) {
	t.Run("should panic for invalid tokens", func(t *testing.T) {
		tests := []keyboard.Token{0, 64, 67}
		for _, token := range tests {
			testName := fmt.Sprintf("for %v", token)
			t.Run(testName, func(t *testing.T) {
				assert.Panics(t, func() {
					keyboard.NewKey(token)
				})
			})
		}
	})
	t.Run("should create key using package variable", func(t *testing.T) {
		key := keyboard.A
		assert.False(t, key.IsUnknown())
		assert.Equal(t, 'A', key.Token().Rune())
	})
	t.Run("should create new key using token", func(t *testing.T) {
		tests := map[string]struct {
			token keyboard.Token
		}{
			"A": {
				token: keyboard.A.Token(),
			},
			"B": {
				token: keyboard.B.Token(),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				key := keyboard.NewKey(test.token)
				assert.Equal(t, test.token, key.Token())
				assert.False(t, key.IsUnknown())
			})
		}
	})
	t.Run("two keys with token should be equal", func(t *testing.T) {
		key1 := keyboard.NewKey(65)
		key2 := keyboard.NewKey(65)
		assert.Equal(t, key1, key2)
	})
	t.Run("two keys with scanCode should be equal", func(t *testing.T) {
		key1 := keyboard.NewUnknownKey(0)
		key2 := keyboard.NewUnknownKey(0)
		assert.Equal(t, key1, key2)
	})
}

func TestNewUnknownKey(t *testing.T) {
	t.Run("should create unknown key using scan code", func(t *testing.T) {
		// when
		key := keyboard.NewUnknownKey(0)
		// then
		assert.True(t, key.IsUnknown())
		assert.Equal(t, 0, key.ScanCode())
	})
}

func TestKeyboard_Pressed(t *testing.T) {
	t.Run("before Update was called, Pressed returns false for all keys", func(t *testing.T) {
		tests := []keyboard.Key{keyboard.A, keyboard.B}
		for _, key := range tests {
			testName := fmt.Sprintf("for key: %v", key)
			t.Run(testName, func(t *testing.T) {
				var (
					event  = keyboard.NewPressedEvent(keyboard.A)
					source = newFakeEventSource(event)
					keys   = keyboard.New(source)
				)
				// when
				pressed := keys.Pressed(key)
				// then
				assert.False(t, pressed)
			})
		}
	})
	t.Run("after Update was called", func(t *testing.T) {
		var (
			aPressed         = keyboard.NewPressedEvent(keyboard.A)
			aReleased        = keyboard.NewReleasedEvent(keyboard.A)
			bPressed         = keyboard.NewPressedEvent(keyboard.B)
			bReleased        = keyboard.NewReleasedEvent(keyboard.B)
			unknown0         = keyboard.NewUnknownKey(0)
			unknown1         = keyboard.NewUnknownKey(1)
			unknown0Pressed  = keyboard.NewPressedEvent(unknown0)
			unknown0Released = keyboard.NewReleasedEvent(unknown0)
			unknown1Released = keyboard.NewReleasedEvent(unknown1)
		)
		tests := map[string]struct {
			source             keyboard.EventSource
			expectedPressed    []keyboard.Key
			expectedNotPressed []keyboard.Key
		}{
			"one PressedEvent for A": {
				source:             newFakeEventSource(aPressed),
				expectedPressed:    []keyboard.Key{keyboard.A},
				expectedNotPressed: []keyboard.Key{keyboard.B},
			},
			"two PressedEvents for B and A": {
				source:          newFakeEventSource(bPressed, aPressed),
				expectedPressed: []keyboard.Key{keyboard.A, keyboard.B},
			},
			"two PressedEvents for A and B": {
				source:          newFakeEventSource(aPressed, bPressed),
				expectedPressed: []keyboard.Key{keyboard.A, keyboard.B},
			},
			"one PressedEvent for unknown key": {
				source:             newFakeEventSource(unknown0Pressed),
				expectedPressed:    []keyboard.Key{unknown0},
				expectedNotPressed: []keyboard.Key{unknown1},
			},
			"two PressedEvents: unknown key and A": {
				source:             newFakeEventSource(unknown0Pressed, aPressed),
				expectedPressed:    []keyboard.Key{unknown0, keyboard.A},
				expectedNotPressed: []keyboard.Key{unknown1},
			},
			"one PressedEvent; one ReleasedEvent for A": {
				source:             newFakeEventSource(aPressed, aReleased),
				expectedNotPressed: []keyboard.Key{keyboard.A},
			},
			"one PressedEvent; one ReleasedEvent for unknown key": {
				source:             newFakeEventSource(unknown0Pressed, unknown0Released),
				expectedNotPressed: []keyboard.Key{unknown0},
			},
			"one PressedEvent for A; one ReleasedEvent for B": {
				source:             newFakeEventSource(aPressed, bReleased),
				expectedPressed:    []keyboard.Key{keyboard.A},
				expectedNotPressed: []keyboard.Key{keyboard.B},
			},
			"one PressedEvent for unknown 0; one ReleasedEvent for unknown 1": {
				source:             newFakeEventSource(unknown0Pressed, unknown1Released),
				expectedPressed:    []keyboard.Key{unknown0},
				expectedNotPressed: []keyboard.Key{unknown1},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				keys := keyboard.New(test.source)
				// when
				keys.Update()
				// then
				for _, expectedPressedKey := range test.expectedPressed {
					assert.True(t, keys.Pressed(expectedPressedKey))
				}
				for _, expectedNotPressedKey := range test.expectedNotPressed {
					assert.False(t, keys.Pressed(expectedNotPressedKey))
				}
			})
		}
	})
}

func newFakeEventSource(events ...keyboard.Event) *fakeEventSource {
	source := &fakeEventSource{}
	source.events = []keyboard.Event{}
	for _, event := range events {
		source.events = append(source.events, event)
	}
	return source
}

type fakeEventSource struct {
	events []keyboard.Event
}

func (f *fakeEventSource) Poll() (keyboard.Event, bool) {
	if len(f.events) > 0 {
		event := f.events[0]
		f.events = f.events[1:]
		return event, true
	}
	return keyboard.EmptyEvent, false
}