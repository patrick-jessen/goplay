package renderer

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/log"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/window"
)

type renderer interface {
	initialize()
	deinitialize()
	render(scene.Scene)
}

var rendererInst renderer

var Settings = settings{
	curType: Forward,
	newType: Forward,
	curAA:   MSAAx4,
	newAA:   MSAAx4,
	curSR:   1024,
	newSR:   1024,
}

type settings struct {
	curType, newType Type
	curAA, newAA     Antialiasing
	curSR, newSR     int
}

type Type int
type Antialiasing int

const (
	Forward Type = iota
)
const (
	NoAA Antialiasing = iota
	FXAA
	MSAAx2
	MSAAx4
	MSAAx8
	MSAAx16
)

func (s *settings) Type() Type {
	return s.curType
}
func (s *settings) SetType(t Type) {
	s.newType = t
}
func (s *settings) Antialiasing() Antialiasing {
	return s.curAA
}
func (s *settings) SetAntialising(a Antialiasing) {
	s.newAA = a

	var maxSamples int32
	gl.GetIntegerv(gl.MAX_SAMPLES, &maxSamples)
	if maxSamples < 16 && s.newAA == MSAAx16 {
		s.newAA = MSAAx8
		log.Warn("MSAA x16 not available, using x8")
	}
	if maxSamples < 8 && s.newAA == MSAAx8 {
		s.newAA = MSAAx4
		log.Warn("MSAA x8 not available, using x4")
	}
	if maxSamples < 4 && s.newAA == MSAAx4 {
		s.newAA = MSAAx2
		log.Warn("MSAA x4 not available, using x2")
	}
	if maxSamples < 2 && s.newAA == MSAAx2 {
		s.newAA = FXAA
		log.Warn("MSAA x2 not available, using FXAA")
	}
}
func (s *settings) ShadowResolution() int {
	return s.curSR
}
func (s *settings) SetShadowResolution(r int) {
	s.newSR = r
}
func (s *settings) Apply() {
	rendererInst.deinitialize()

	if s.curType != s.newType {
		switch s.newType {
		case Forward:
			rendererInst = &forwardRenderer{}
		default:
			log.Panic("not implemented")
		}
	}

	s.curType = s.newType
	s.curAA = s.newAA
	s.curSR = s.newSR

	rendererInst.initialize()
}

func onResize(w, h int) {
	rendererInst.deinitialize()
	rendererInst.initialize()
}
func Initialize() {
	rendererInst = &forwardRenderer{}
	rendererInst.initialize()

	window.AddResizeHandler(onResize)
}
func Deinitialize() {
	rendererInst.deinitialize()
}
func Render() {
	rendererInst.render(scene.Current())
}
