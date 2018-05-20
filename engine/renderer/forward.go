package renderer

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/framebuffer"
	"github.com/patrick-jessen/goplay/engine/model"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/shader"
	"github.com/patrick-jessen/goplay/engine/window"
)

type quadMat struct {
	Shader shader.Shader
}

func (m quadMat) Apply() {
	m.Shader.Use()
}

type forwardRenderer struct {
	shaderFrameBuffer *framebuffer.FrameBuffer
	width, height     int
	postScene         scene.Scene
}

func (f *forwardRenderer) initialize() {
	f.postScene = scene.New()

	f.width, f.height = window.Settings.Size()

	var msLevel int
	switch Settings.curAA {
	case NoAA:
	case FXAA:
		model.Load("quad").Mount(f.postScene.Root)
		s := shader.Load("fxaa")
		uRes := s.GetUniform("resolution")
		gl.Uniform2f(uRes, float32(f.width), float32(f.height))

		f.postScene.Root.Child("0").Component("MeshRenderer").(*model.MeshRenderer).Mat = &quadMat{
			Shader: s,
		}
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

func (f *forwardRenderer) deinitialize() {
	f.shaderFrameBuffer.Free()
}

func (f *forwardRenderer) render(scene scene.Scene) {

	// Shadow map pass
	f.renderShadows()

	// Shading pass
	f.shaderFrameBuffer.Bind()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	scene.Render()

	// Postprocessing pass
	switch Settings.curAA {
	case FXAA:
		framebuffer.Unbind()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		f.shaderFrameBuffer.BindColorTexture(0, 0)
		f.postScene.Render()
	default:
		// 1. NoAA blits onto default frame buffer 1:1.
		// 2. MSAAx_ blits onto default frame buffer and
		// performs linear interpolation on samples.
		f.shaderFrameBuffer.Blit(nil, f.width, f.height, false)
	}
}

func (f *forwardRenderer) renderShadows() {
	// TODO
}
