package window

import (
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

var windowSettings = struct {
	width, height int
	fullscreen    bool
	title         string
	vsync         bool
	msaa          int
}{
	width:      800,
	height:     600,
	title:      "GoPlay",
	vsync:      false,
	msaa:       4,
	fullscreen: false,
}

var winHandle *glfw.Window
var resizeHandlers []func(int, int)

// init initializes the window system.
func init() {
	if err := glfw.Init(); err != nil {
		panic("failed to initialize window system:\n" + err.Error())
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
}

// Deinitialize deinitializes the window system.
// Should be called before exiting application.
func Deinitialize() {
	glfw.Terminate()
}

// Create creates the window.
func Create() {
	var err error

	// mode := glfw.GetPrimaryMonitor().GetVideoMode()
	// glfw.WindowHint(glfw.RedBits, mode.RedBits)
	// glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
	// glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
	// glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)
	glfw.WindowHint(glfw.Samples, windowSettings.msaa)

	var fsMon *glfw.Monitor
	if windowSettings.fullscreen {
		fsMon = glfw.GetPrimaryMonitor()
	}
	winHandle, err = glfw.CreateWindow(windowSettings.width, windowSettings.height, windowSettings.title, fsMon, nil)
	if err != nil {
		panic("failed to create window:\n" + err.Error())
	}
	winHandle.MakeContextCurrent()
	winHandle.SetKeyCallback(keyCallback)
	winHandle.SetMouseButtonCallback(mouseButtonCallback)
	winHandle.SetCursorPosCallback(cursorPosCallback)
	winHandle.SetScrollCallback(scrollCallback)
	winHandle.SetFramebufferSizeCallback(resizeCallback)

	if err := gl.Init(); err != nil {
		panic("failed to initialize OpenGL:/n" + err.Error())
	}

	SetVerticalSync(windowSettings.vsync)

	gl.Viewport(0, 0, int32(windowSettings.width), int32(windowSettings.height))
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.FRAMEBUFFER_SRGB)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)
}

// AddResizeHandler sets the resize handler.
func AddResizeHandler(handler func(int, int)) {
	resizeHandlers = append(resizeHandlers, handler)
}

// Close closes the window.
func Close() {
	winHandle.Destroy()
	glfw.Terminate()
}

// SetTitle sets the window title.
func SetTitle(title string) {
	if winHandle != nil {
		winHandle.SetTitle(title)
	}
	windowSettings.title = title
}

// SetVideoMode set the video mode for the window.
// if fs is true, the window will become full-screen.
func SetVideoMode(fs bool, w, h int) {
	vm := glfw.GetPrimaryMonitor().GetVideoMode()
	refresh := vm.RefreshRate

	if fs {
		if w == -1 || h == -1 {
			w = vm.Width
			h = vm.Height
		}
		if winHandle != nil {
			winHandle.SetMonitor(glfw.GetPrimaryMonitor(), 0, 0, w, h, refresh)
			SetVerticalSync(windowSettings.vsync)
		}
	} else {
		if w == -1 || h == -1 {
			w = 800
			h = 600
		}
		if winHandle != nil {
			winHandle.SetMonitor(nil, 100, 100, w, h, refresh)
			SetVerticalSync(windowSettings.vsync)
		}
	}

	windowSettings.fullscreen = fs
	windowSettings.width = w
	windowSettings.height = h
}

// SetVerticalSync enables or disables vertical synchronization.
func SetVerticalSync(on bool) {
	windowSettings.vsync = on
	if winHandle != nil {
		interval := 0
		if on {
			interval = 1
		}
		glfw.SwapInterval(interval)
	}
}

// VerticalSync return whether vsync is enabled or not.
func VerticalSync() bool {
	return windowSettings.vsync
}

// Size returns the window size.
func Size() [2]int {
	return [2]int{windowSettings.width, windowSettings.height}
}

// ShouldClose indicates whether the window should close.
// Should be used to create the main loop.
// Eg.
// for window.ShouldClose() {
//   // do stuff
//   window.Update()
// }
func ShouldClose() bool {
	if winHandle == nil {
		return true
	}
	return winHandle.ShouldClose()
}

// Update polls events and swaps the content to front.
// For usage, see ShouldClose().
func Update() {
	updateInput()
	glfw.PollEvents()
	winHandle.SwapBuffers()
}

// resizeCallback handles window resize.
func resizeCallback(w *glfw.Window, width, height int) {
	for _, r := range resizeHandlers {
		r(width, height)
	}
	windowSettings.width = width
	windowSettings.height = height
	gl.Viewport(0, 0, int32(windowSettings.width), int32(windowSettings.height))
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
