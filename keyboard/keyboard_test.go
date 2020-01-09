package keyboard_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func Test(t *testing.T) {
	t.Run("should create key using package variable", func(t *testing.T) {
		key := keyboard.A
		assert.False(t, key.IsUnknown())
		assert.Equal(t, keyboard.Token("A"), key.Token())
	})
	t.Run("two keys with same token should be equal", func(t *testing.T) {
		key1 := keyboard.A
		key2 := keyboard.A
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

func TestPressedKeys(t *testing.T) {
	var (
		aPressed         = keyboard.NewPressedEvent(keyboard.A)
		aReleased        = keyboard.NewReleasedEvent(keyboard.A)
		bPressed         = keyboard.NewPressedEvent(keyboard.B)
		unknown0         = keyboard.NewUnknownKey(0)
		unknown1         = keyboard.NewUnknownKey(0)
		unknown0Pressed  = keyboard.NewPressedEvent(unknown0)
		unknown0Released = keyboard.NewReleasedEvent(unknown0)
		unknown1Pressed  = keyboard.NewPressedEvent(unknown1)
	)
	t.Run("before Update pressed keys are empty", func(t *testing.T) {
		source := newFakeEventSource(aPressed)
		keys := keyboard.New(source)
		// when
		pressed := keys.PressedKeys()
		// then
		assert.Empty(t, pressed)
	})
	t.Run("after Update", func(t *testing.T) {
		tests := map[string]struct {
			source          keyboard.EventSource
			expectedPressed []keyboard.Key
		}{
			"one PressedEvent for A": {
				source:          newFakeEventSource(aPressed),
				expectedPressed: []keyboard.Key{keyboard.A},
			},
			"one PressedEvent for B": {
				source:          newFakeEventSource(bPressed),
				expectedPressed: []keyboard.Key{keyboard.B},
			},
			"one PressedEvent for A, one ReleaseEvent for A": {
				source:          newFakeEventSource(aPressed, aReleased),
				expectedPressed: nil,
			},
			"one PressedEvent for unknown 0": {
				source:          newFakeEventSource(unknown0Pressed),
				expectedPressed: []keyboard.Key{unknown0},
			},
			"one PressedEvent for unknown 1": {
				source:          newFakeEventSource(unknown1Pressed),
				expectedPressed: []keyboard.Key{unknown1},
			},
			"one PressedEvent for unknown 0, one ReleaseEvent for unknown 0": {
				source:          newFakeEventSource(unknown0Pressed, unknown0Released),
				expectedPressed: nil,
			},
			"A pressed, then released and pressed again": {
				source:          newFakeEventSource(aPressed, aReleased, aPressed),
				expectedPressed: []keyboard.Key{keyboard.A},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				keys := keyboard.New(test.source)
				keys.Update()
				// when
				pressed := keys.PressedKeys()
				// then
				assert.Equal(t, test.expectedPressed, pressed)
			})
		}
	})
	t.Run("after second update", func(t *testing.T) {
		t.Run("when A was pressed before first update", func(t *testing.T) {
			source := newFakeEventSource(aPressed)
			keys := keyboard.New(source)
			keys.Update()
			keys.Update()
			// when
			pressed := keys.PressedKeys()
			// then
			assert.Equal(t, []keyboard.Key{keyboard.A}, pressed)
		})
		t.Run("when A was pressed before first update, then released before second one", func(t *testing.T) {
			source := newFakeEventSource(aPressed)
			keys := keyboard.New(source)
			keys.Update()
			source.events = append(source.events, aReleased)
			keys.Update()
			// when
			pressed := keys.PressedKeys()
			// then
			assert.Empty(t, pressed)
		})
	})
}

func TestJustPressed(t *testing.T) {
	var (
		aPressed        = keyboard.NewPressedEvent(keyboard.A)
		aReleased       = keyboard.NewReleasedEvent(keyboard.A)
		bPressed        = keyboard.NewPressedEvent(keyboard.B)
		unknown0Pressed = keyboard.NewPressedEvent(keyboard.NewUnknownKey(0))
	)

	t.Run("before update should return false", func(t *testing.T) {
		tests := map[string]struct {
			source keyboard.EventSource
		}{
			"for no events": {
				source: newFakeEventSource(),
			},
			"when A has been pressed": {
				source: newFakeEventSource(aPressed),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				keys := keyboard.New(test.source)
				// when
				justPressed := keys.JustPressed(keyboard.A)
				// then
				assert.False(t, justPressed)
			})
		}
	})

	t.Run("after first update", func(t *testing.T) {
		tests := map[string]struct {
			source              keyboard.EventSource
			expectedJustPressed bool
		}{
			"for no events": {
				source:              newFakeEventSource(),
				expectedJustPressed: false,
			},
			"when A has been pressed": {
				source:              newFakeEventSource(aPressed),
				expectedJustPressed: true,
			},
			"when B has been pressed": {
				source:              newFakeEventSource(bPressed),
				expectedJustPressed: false,
			},
			"when Unknown 0 has been pressed": {
				source:              newFakeEventSource(unknown0Pressed),
				expectedJustPressed: false,
			},
			"when A has been released": {
				source:              newFakeEventSource(aReleased),
				expectedJustPressed: false,
			},
			"when A has been pressed and released": {
				source:              newFakeEventSource(aPressed, aReleased),
				expectedJustPressed: true,
			},
			"when A has been released and pressed": {
				source:              newFakeEventSource(aReleased, aPressed),
				expectedJustPressed: true,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				keys := keyboard.New(test.source)
				keys.Update()
				// when
				justPressed := keys.JustPressed(keyboard.A)
				// then
				assert.Equal(t, test.expectedJustPressed, justPressed)
			})
		}
	})

	t.Run("should return false after second update", func(t *testing.T) {
		source := newFakeEventSource(aPressed)
		keys := keyboard.New(source)
		keys.Update()
		keys.Update()
		// when
		pressed := keys.JustPressed(keyboard.A)
		// then
		assert.False(t, pressed)
	})
}

func TestKey_Serialize(t *testing.T) {
	t.Run("should serialize key", func(t *testing.T) {
		tests := map[string]struct {
			key                      keyboard.Key
			expectedSerializedString string
		}{
			"A": {
				key:                      keyboard.A,
				expectedSerializedString: "A",
			},
			"B": {
				key:                      keyboard.B,
				expectedSerializedString: "B",
			},
			"Home": {
				key:                      keyboard.Home,
				expectedSerializedString: "Home",
			},
			"KeypadAdd": {
				key:                      keyboard.KeypadAdd,
				expectedSerializedString: "Keypad +",
			},
			"Unknown 0": {
				key:                      keyboard.NewUnknownKey(0),
				expectedSerializedString: "?0",
			},
			"Unknown 1": {
				key:                      keyboard.NewUnknownKey(1),
				expectedSerializedString: "?1",
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				serialized := test.key.Serialize()
				// then
				assert.Equal(t, test.expectedSerializedString, serialized)
			})
		}
	})
}

func TestKey_Deserialize(t *testing.T) {
	t.Run("should return error", func(t *testing.T) {
		tests := map[string]string{
			"<empty>": "",
			"!":       "!",
			"#":       "#",
			"$":       "$",
			"?":       "?",
			"?A":      "?A",
		}
		for name, serializedString := range tests {
			t.Run(name, func(t *testing.T) {
				_, err := keyboard.Deserialize(serializedString)
				assert.Error(t, err)
			})
		}
	})
	t.Run("should deserialize key", func(t *testing.T) {
		tests := map[string]struct {
			serializedString string
			expectedKey      keyboard.Key
		}{
			"A": {
				serializedString: "A",
				expectedKey:      keyboard.A,
			},
			"B": {
				serializedString: "B",
				expectedKey:      keyboard.B,
			},
			"Home": {
				serializedString: "Home",
				expectedKey:      keyboard.Home,
			},
			"KeypadAdd": {
				serializedString: "Keypad +",
				expectedKey:      keyboard.KeypadAdd,
			},
			"?0": {
				serializedString: "?0",
				expectedKey:      keyboard.NewUnknownKey(0),
			},
			"?1": {
				serializedString: "?1",
				expectedKey:      keyboard.NewUnknownKey(1),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				key, err := keyboard.Deserialize(test.serializedString)
				// then
				require.NoError(t, err)
				assert.Equal(t, test.expectedKey, key)
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
