package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/framebuffer"
	"github.com/patrick-jessen/goplay/engine/model"
	"github.com/patrick-jessen/goplay/engine/renderer"

	"github.com/patrick-jessen/goplay/components"
	"github.com/patrick-jessen/goplay/engine"
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

type app struct {
	scene     scene.Scene
	camera    *components.Camera
	postScene scene.Scene

	fb *framebuffer.FrameBuffer
}

func (a *app) OnStart() {
	gl.ClearColor(1, 1, 1, 1)
	a.postScene = scene.New()
	model.Load("quad").Mount(a.postScene.Root)

	s := shader.Load("fxaa")
	uRes := s.GetUniform("resolution")
	w, h := window.Settings.Size()
	gl.Uniform2f(uRes, float32(w), float32(h))

	a.postScene.Root.Child("0").Component("MeshRenderer").(*model.MeshRenderer).Mat = &quadMat{
		Shader: s,
	}

	a.scene = scene.Load("main")
	a.scene.MakeCurrent()
	a.camera = a.scene.Root.Child("camera").Component("Camera").(*components.Camera)

	for k, v := range scene.MountMap {
		model.Load(v).Mount(k)
	}

	a.fb = framebuffer.New(w, h, 1, 0)
}

func (a *app) OnUpdate() {

	a.scene.Update() // <- move to engine
	shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())

	renderer.Render(a.scene)
	return

	// FXAA
	a.fb.Bind()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())
	a.scene.Render()

	// Render post-processing quad
	framebuffer.Unbind()
	a.fb.BindColorTexture(0, 0)
	a.fb.BindDepthTexture(1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	a.postScene.Render()

}
func (a *app) OnExit() {}

func main() {
	window.Settings.SetVSync(true)
	window.Settings.SetTitle("MyGame")
	window.Settings.SetSize(1024, 768)

	engine.Start(&app{})
}
