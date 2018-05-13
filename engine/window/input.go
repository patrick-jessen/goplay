package window

import mgl "github.com/go-gl/mathgl/mgl64"

// MouseMoveInput is an event for when the mouse is moved.
type MouseMoveInput struct {
	X, Y int
}

// MouseScrollInput is an event for when the mouse wheel is scrolled.
type MouseScrollInput struct {
	Delta float64
}

// MouseButtonInput is an event for when a mouse button is pressed.
type MouseButtonInput struct {
	Button  int
	Press   bool
	Release bool
}

// KeyboardInput is an event for when a key is pressed.
type KeyboardInput struct {
	Key     int
	Press   bool
	Release bool
}

var (
	lastMousePosition  mgl.Vec2
	mousePosition      mgl.Vec2
	mouseScroll        float64
	mouseButtonPressed [3]bool
	mouseButtonEvent   [3]int
)

// MousePosition returns the current mouse position.
func MousePosition() mgl.Vec2 {
	return mousePosition
}

// MouseMove returns the delta that the mouse moved since last frame.
func MouseMove() mgl.Vec2 {
	return mousePosition.Sub(lastMousePosition)
}

// MouseScroll returns the delta that the mouse wheel scrolled since last frame.
func MouseScroll() float64 {
	return mouseScroll
}

// MouseButton returns whether a mouse button was pressed since last frame.
func MouseButton(b int) bool {
	return mouseButtonPressed[b]
}

// MouseButtonDown returns whether a mouse button is currently down.
func MouseButtonDown(b int) bool {
	return mouseButtonEvent[b] == 1
}

// MouseButtonUp returns whether a mouse button is currently up.
func MouseButtonUp(b int) bool {
	return mouseButtonEvent[b] == -1
}

// updateInput is called each main loop to reset input.
func updateInput() {
	lastMousePosition = mousePosition
	mouseScroll = 0
	mouseButtonEvent = [3]int{}
}

// onInput is called by the window when an input event occurs.
func onInput(e interface{}) {
	switch e.(type) {
	case *MouseMoveInput:
		tmp := e.(*MouseMoveInput)
		mousePosition = mgl.Vec2{float64(tmp.X), float64(tmp.Y)}

	case *MouseScrollInput:
		tmp := e.(*MouseScrollInput)
		mouseScroll = tmp.Delta

	case *MouseButtonInput:
		tmp := e.(*MouseButtonInput)
		if tmp.Press {
			mouseButtonPressed[tmp.Button] = true
			mouseButtonEvent[tmp.Button] = 1
		} else if tmp.Release {
			mouseButtonPressed[tmp.Button] = false
			mouseButtonEvent[tmp.Button] = -1
		}

		// case *KeyboardInput:
		// 	keyboard = e.(*KeyboardInput)
	}
}
