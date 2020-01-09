// Package keyboard adds support for keyboard input.
//
// You can start using using keyboard by creating Keyboard instance:
//
//     window := windows.Open(...)
//     keys := keyboard.New(window)
//     loops.Loop(window, func(frame *pixiq.Frame) {
//         keys.Update() // This is needed each frame to update the state of keys
//         if keys.Pressed(keyboard.A) {
//             ...
//         }
//     })
//
package keyboard

import (
	"fmt"
	"strconv"
	"strings"
)

// EventSource is a source of keyboard Events.
type EventSource interface {
	// Poll retrieves and removes next keyboard Event. If there are no more
	// events false is returned.
	Poll() (Event, bool)
}

func newKey(token Token) Key {
	return Key{token: token}
}

// Token is a unique string representation of the key. It the key is
// unknown then token is empty.
type Token string

// NewUnknownKey creates key with platform-specific scanCode.
func NewUnknownKey(scanCode int) Key {
	return Key{
		scanCode: scanCode,
	}
}

// Key identifies the pressed or released key.
type Key struct {
	token    Token
	scanCode int
}

// IsUnknown returns true if mapping has not been found. ScanCode should be used
// instead.
func (k Key) IsUnknown() bool {
	return k.token == ""
}

// ScanCode returns the platform-specific code.
func (k Key) ScanCode() int {
	return k.scanCode
}

// Token returns a platform-independent unique token.
func (k Key) Token() Token {
	return k.token
}

// Serialize marshals the key to string which can be used for saving action keymap.
func (k Key) Serialize() string {
	if k.IsUnknown() {
		return fmt.Sprintf("?%d", k.scanCode)
	}
	return string(k.Token())
}

// Deserialize unmarshalls the key from string which can be used for loading action keymap.
func Deserialize(s string) (Key, error) {
	if strings.HasPrefix(s, "?") && len(s) > 1 {
		scanCode, err := strconv.Atoi(s[1:])
		if err != nil {
			return Key{}, fmt.Errorf("unserializable key string %s: %s", s, err)
		}
		return NewUnknownKey(scanCode), nil
	}
	var found bool
	for _, key := range allKeys {
		if string(key.token) == s {
			found = true
			break
		}
	}
	if !found {
		return Key{}, fmt.Errorf("unserializable key string %s", s)
	}
	return newKey(Token(s)), nil
}

// EmptyEvent should be returned by EventSource when it does not have more events.
var EmptyEvent = Event{}

// NewPressedEvent returns new instance of Event when key was pressed.
func NewPressedEvent(key Key) Event {
	return Event{
		typ: pressed,
		key: key,
	}
}

// NewReleasedEvent returns new instance of Event when key was released.
func NewReleasedEvent(key Key) Event {
	return Event{
		typ: released,
		key: key,
	}
}

// Event describes what happened with the key. Whether it was pressed or released.
type Event struct {
	typ eventType
	key Key
}

// eventType is used because using polymorphism means heap allocation and we don't
// want to generate garbage.
type eventType byte

const (
	pressed  eventType = 1
	released eventType = 2
)

// New creates Keyboard instance.
func New(source EventSource) *Keyboard {
	if source == nil {
		panic("nil EventSource")
	}
	return &Keyboard{
		source:      source,
		pressed:     make(map[Key]struct{}),
		justPressed: make(map[Key]bool),
	}
}

// Keyboard provides a read-only information about the current state of the
// keyboard, such as what keys are currently pressed. Please note that
// updating the Keyboard state retrieves and removes events from EventSource.
// Therefore only Keyboard instance can be created for one EventSource.
type Keyboard struct {
	source      EventSource
	pressed     map[Key]struct{}
	justPressed map[Key]bool
}

// Update updates the state of the keyboard by polling events queued since last
// time the function was executed.
func (k *Keyboard) Update() {
	k.clearJustPressed()
	for {
		event, ok := k.source.Poll()
		if !ok {
			return
		}
		switch event.typ {
		case pressed:
			k.pressed[event.key] = struct{}{}
			k.justPressed[event.key] = true
		case released:
			delete(k.pressed, event.key)
		}
	}
}

func (k *Keyboard) clearJustPressed() {
	for key := range k.justPressed {
		delete(k.justPressed, key)
	}
}

// Pressed returns true if given key is currently pressed.
// If during two last keyboard.Update calls the key was pressed and released
// then the this method returns false.
func (k *Keyboard) Pressed(key Key) bool {
	_, found := k.pressed[key]
	return found
}

// PressedKeys returns a slice of all currently pressed keys. It may be empty
// aka nil. This function can be used to get a key mapping for a given action
// in the game.
// If during two last keyboard.Update calls the key was pressed and released
// then the key is not returned.
func (k *Keyboard) PressedKeys() []Key {
	var pressedKeys []Key
	for key := range k.pressed {
		pressedKeys = append(pressedKeys, key)
	}
	return pressedKeys
}

// JustPressed returns true if the key was pressed between two last keyboard.Update
// calls. If it was pressed and released at the same time between these calls
// this method return true.
func (k *Keyboard) JustPressed(key Key) bool {
	return k.justPressed[key]
}
