package window

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/patrick-jessen/goplay/engine/log"
)

// Settings is the window settings.
var Settings settings

type settings struct {
	curVSync bool
	newVSync bool

	curTitle string
	newTitle string

	newSize [2]int
	newFS   bool
}

func (s *settings) SetVSync(on bool) {
	s.newVSync = on
}
func (s *settings) VSync() bool {
	return s.curVSync
}

func (s *settings) SetTitle(str string) {
	s.newTitle = str
}
func (s *settings) Title() string {
	return s.curTitle
}

func (s *settings) SetSize(width, height int) {
	s.newSize = [2]int{width, height}
}
func (s *settings) Size() (width, height int) {
	return winHandle.GetSize()
}

func (s *settings) SetFullScreen(on bool) {
	s.newFS = on
}
func (s *settings) FullScreen() bool {
	return winHandle.GetMonitor() != nil
}

func (s *settings) Apply() {
	if winHandle == nil {
		log.Warn("cannot apply settings before window is created")
		return
	}

	// Apply title
	if s.Title() != s.newTitle {
		winHandle.SetTitle(s.newTitle)
		s.curTitle = s.newTitle
	}

	// Apply size
	if w, h := s.Size(); w != s.newSize[0] || h != s.newSize[1] {
		winHandle.SetSize(s.newSize[0], s.newSize[1])
	}

	// Apply full screen
	if s.FullScreen() != s.newFS {
		var mon *glfw.Monitor
		if s.newFS {
			mon = glfw.GetPrimaryMonitor()
		}
		winHandle.SetMonitor(mon, 0, 0, s.newSize[0], s.newSize[1], glfw.DontCare)

		// Force reset of vsync
		s.curVSync = !s.newVSync
	}

	// Apply vsync
	if s.VSync() != s.newVSync {
		interval := 0
		if s.newVSync {
			interval = 1
		}
		glfw.SwapInterval(interval)
		s.curVSync = s.newVSync
	}
}
