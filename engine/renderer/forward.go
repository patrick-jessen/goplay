package renderer

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/framebuffer"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/window"
)

type forwardRenderer struct {
	scene scene.Scene

	shaderFrameBuffer *framebuffer.FrameBuffer
	width, height     int
}

func NewForward() *forwardRenderer {
	return &forwardRenderer{}
}

func (f *forwardRenderer) Initialize(s scene.Scene) {
	f.scene = s

	// Locate all lights

	// Subscribe to AddNode()
	// so that future ligts can be captured

	f.width, f.height = window.Settings.Size()

	var msLevel int
	switch Settings.curAA {
	case NoAA:
	case FXAA:
	case MSAAx2:
		msLevel = 2
	case MSAAx4:
		msLevel = 4
	case MSAAx8:
		msLevel = 8
	case MSAAx16:
		msLevel = 16
	}

	f.shaderFrameBuffer = framebuffer.New(f.width, f.height, 1, msLevel)
}

func (f *forwardRenderer) Render() {

	// Shadow map pass
	f.renderShadows()

	// Shading pass
	f.shaderFrameBuffer.Bind()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	f.scene.Render()

	// Postprocessing pass
	switch Settings.curAA {
	case NoAA: // TODO
	case FXAA: // TODO
		// Draw quad
	default:
		f.shaderFrameBuffer.Blit(nil, f.width, f.height, false)
	}
}

func (f *forwardRenderer) renderShadows() {

}
