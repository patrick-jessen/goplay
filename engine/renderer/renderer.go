package renderer

import (
	"github.com/patrick-jessen/goplay/engine/scene"
)

var Settings = settings{
	curAA: MSAAx4,
	newAA: MSAAx4,
}

type Antialiasing int

const (
	NoAA Antialiasing = iota
	FXAA
	MSAAx2
	MSAAx4
	MSAAx8
	MSAAx16
)

type ShadowQuality int

const (
	ShadowOff    ShadowQuality = 0
	ShadowLow    ShadowQuality = 512
	ShadowMedium ShadowQuality = 1024
	ShadowHigh   ShadowQuality = 2048
)

type settings struct {
	curAA, newAA Antialiasing
	curSQ, newSQ ShadowQuality
}

func (s *settings) Antialiasing() Antialiasing {
	return s.curAA
}
func (s *settings) SetAntialising(a Antialiasing) {
	s.newAA = a
}

type Renderer interface {
	Initialize(scene.Scene)
	Render()
}
