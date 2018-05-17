package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/framebuffer"
	"github.com/patrick-jessen/goplay/engine/model"

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
	a.postScene = scene.New()
	model.Load("quad").Mount(a.postScene.Root)
	a.postScene.Root.Child("0").Component("MeshRenderer").(*model.MeshRenderer).Mat = &quadMat{
		Shader: shader.Load("quad"),
	}

	a.scene = scene.Load("main")
	a.scene.MakeCurrent()
	a.camera = a.scene.Root.Child("camera").Component("Camera").(*components.Camera)

	for k, v := range scene.MountMap {
		model.Load(v).Mount(k)
	}

	a.fb = framebuffer.New(1)
}
func (a *app) OnUpdate() {

	a.scene.Update() // <- move to engine

	// Render shadow maps
	// TODO

	// Render normally
	framebuffer.Use(a.fb)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())
	a.scene.Render()

	// Render post-processing quad

	framebuffer.Use(nil)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	a.fb.BindTexture(0, 0)
	a.postScene.Render()

}
func (a *app) OnExit() {}

func main() {
	window.Settings.SetVSync(true)
	window.Settings.SetTitle("MyGame")
	window.Settings.SetSize(1024, 768)

	engine.Start(&app{})
}
