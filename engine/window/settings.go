package window

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/patrick-jessen/goplay/engine/log"
)

var Settings settings

type settings struct {
	curVSync bool
	newVSync bool

	curTitle string
	newTitle string

	curRes [2]int
	newRes [2]int
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
func (s *settings) SetResolution(width, height int) {
	s.newRes = [2]int{width, height}
}
func (s *settings) Resolution() (width, height int) {
	return s.curRes[0], s.curRes[1]
}

func (s *settings) Apply() {
	if winHandle == nil {
		log.Warn("cannot apply settings before window is created")
		return
	}

	// Apply vsync
	if s.curVSync != s.newVSync {
		interval := 0
		if s.newVSync {
			interval = 1
		}
		glfw.SwapInterval(interval)
		s.curVSync = s.newVSync
	}

	// Apply title
	if s.curTitle != s.newTitle {
		winHandle.SetTitle(s.newTitle)
		s.curTitle = s.newTitle
	}
}
