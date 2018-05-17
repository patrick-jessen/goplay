package window

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

var winHandle *glfw.Window
var resizeHandlers []func(int, int)

// Create creates the window.
func Create() {
	if err := glfw.Init(); err != nil {
		log.Panic("failed to initialize window system", "error", err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	mode := glfw.GetPrimaryMonitor().GetVideoMode()
	glfw.WindowHint(glfw.RedBits, mode.RedBits)
	glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
	glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
	glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)

	var err error
	winHandle, err = glfw.CreateWindow(800, 600, "GoPlay", nil, nil)
	if err != nil {
		log.Panic("failed to create window", "error", err)
	}
	winHandle.MakeContextCurrent()

	winHandle.SetKeyCallback(keyCallback)
	winHandle.SetMouseButtonCallback(mouseButtonCallback)
	winHandle.SetCursorPosCallback(cursorPosCallback)
	winHandle.SetScrollCallback(scrollCallback)
	winHandle.SetFramebufferSizeCallback(resizeCallback)

	if err := gl.Init(); err != nil {
		log.Panic("failed to initialize OpenGL", "error", err)
	}

	Settings.Apply()

	w, h := Settings.Size()
	gl.Viewport(0, 0, int32(w), int32(h))
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.FRAMEBUFFER_SRGB)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)
}

// Destroy closes the window.
func Destroy() {
	winHandle.Destroy()
	glfw.Terminate()
}

// AddResizeHandler sets the resize handler.
func AddResizeHandler(handler func(int, int)) {
	resizeHandlers = append(resizeHandlers, handler)
}

// ShouldClose indicates whether the window should close.
// Should be used to create the main loop.
// Eg.
// for window.ShouldClose() {
//   // do stuff
//   window.Update()
// }
func ShouldClose() bool {
	return winHandle.ShouldClose()
}

// Update polls events and swaps the content to front.
// For usage, see ShouldClose().
func Update() {
	winHandle.SwapBuffers()
	updateInput()
	glfw.PollEvents()
}

// resizeCallback handles window resize.
func resizeCallback(w *glfw.Window, width, height int) {
	for _, r := range resizeHandlers {
		r(width, height)
	}
	gl.Viewport(0, 0, int32(width), int32(height))
}

// keyCallback is called when a key is pressed.
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	onInput(&KeyboardInput{
		Key:     int(key),
		Press:   action == 1,
		Release: action == 0,
	})
}

// mouseButtonCallback is called when a mouse button is pressed.
func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	onInput(&MouseButtonInput{
		Button:  int(button),
		Press:   action == 1,
		Release: action == 0,
	})
}

// cursorPosCallback is called when the mouse is moved.
func cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	onInput(&MouseMoveInput{int(xpos), int(ypos)})
}

// scrollCallback is called when the mouse wheel is scrolled.
func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	onInput(&MouseScrollInput{float32(yoff)})
}
