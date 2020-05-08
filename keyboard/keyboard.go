// Package keyboard adds support for keyboard input.
//
// You can start using keyboard by creating Keyboard instance:
//
//     keys := keyboard.New(window)
//     loop.Run(window, func(frame *loop.Frame) {
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

// EventSource is a source of keyboard Events. On each Update() Keyboard polls
// the EventSource by executing PollKeyboardEvent method multiple times - until PollKeyboardEvent()
// returns false. In other words Keyboard#Update drains the EventSource.
type EventSource interface {
	// PollKeyboardEvent retrieves and removes next keyboard Event. If there are no more
	// events false is returned.
	PollKeyboardEvent() (Event, bool)
}

func newKey(token token) Key {
	return Key{token: token}
}

// token is a unique string representation of the key. It the key is
// unknown then token is empty.
type token string

// NewUnknownKey creates key with platform-specific scanCode.
func NewUnknownKey(scanCode int) Key {
	return Key{
		scanCode: scanCode,
	}
}

// Key identifies the pressed or released key.
type Key struct {
	token    token
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

// Serialize marshals the key to string which can be used for saving action keymap.
func (k Key) Serialize() string {
	if k.IsUnknown() {
		return fmt.Sprintf("?%d", k.scanCode)
	}
	return string(k.token)
}

// String returns the string representation of the Key for debugging purposes. Do not
// use this method for serialization. Use Key.Serialize instead.
func (k Key) String() string {
	if k.IsUnknown() {
		return fmt.Sprintf("Key with scanCode %d", k.scanCode)
	}
	return "Key " + string(k.token)
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
	for _, key := range AllKeys {
		if string(key.token) == s {
			found = true
			break
		}
	}
	if !found {
		return Key{}, fmt.Errorf("unserializable key string %s", s)
	}
	return newKey(token(s)), nil
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
//
// Event can be constructed using NewXXXEvent function
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

// New creates Keyboard instance. It will consume all events from EventSource each
// time Update method is called. For this reason you can't have two Keyboard instances
// for the same EventSource.
func New(source EventSource) *Keyboard {
	if source == nil {
		panic("nil EventSource")
	}
	return &Keyboard{
		source:       source,
		pressed:      make(map[Key]struct{}),
		justPressed:  make(map[Key]bool),
		justReleased: make(map[Key]bool),
	}
}

// Keyboard provides a read-only information about the current state of the
// keyboard, such as what keys are currently pressed. Please note that
// updating the Keyboard state retrieves and removes events from EventSource.
// Therefore only one Keyboard instance can be created for specific EventSource.
type Keyboard struct {
	source       EventSource
	pressed      map[Key]struct{}
	justPressed  map[Key]bool
	justReleased map[Key]bool
}

// Update updates the state of the keyboard by polling events queued since last
// time the function was executed.
func (k *Keyboard) Update() {
	k.clearJustPressed()
	k.clearJustReleased()
	for {
		event, ok := k.source.PollKeyboardEvent()
		if !ok {
			return
		}
		switch event.typ {
		case pressed:
			k.pressed[event.key] = struct{}{}
			k.justPressed[event.key] = true
		case released:
			delete(k.pressed, event.key)
			k.justReleased[event.key] = true
		}
	}
}

func (k *Keyboard) clearJustPressed() {
	for key := range k.justPressed {
		delete(k.justPressed, key)
	}
}

func (k *Keyboard) clearJustReleased() {
	for key := range k.justReleased {
		delete(k.justReleased, key)
	}
}

// Pressed returns true if given key is currently pressed.
// If between two last keyboard.Update calls the key was pressed and released
// then the this method returns false.
func (k *Keyboard) Pressed(key Key) bool {
	_, found := k.pressed[key]
	return found
}

// PressedKeys returns a slice of all currently pressed keys. It may be empty
// aka nil. This function can be used to get a key mapping for a given action
// in the game.
// If between two last keyboard.Update calls the key was pressed and released
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

// JustReleased returns true if the key was released between two last keyboard.Update
// calls. If it was released and pressed at the same time between these calls
// this method return true.
func (k *Keyboard) JustReleased(key Key) bool {
	return k.justReleased[key]
}
