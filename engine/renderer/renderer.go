package renderer

import (
	"github.com/patrick-jessen/goplay/engine/scene"
)

var Settings = settings{}

type Antialiasing int

const (
	NoAA Antialiasing = iota
	FXAA
	MSAAx2
	MSAAx4
	MSAAx8
	MSAAx16
)

type settings struct {
}

type Renderer interface {
	Render(scene.Scene)
}
