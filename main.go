package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/model"

	"github.com/patrick-jessen/goplay/components"
	"github.com/patrick-jessen/goplay/engine"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/shader"
	"github.com/patrick-jessen/goplay/engine/window"
)

type app struct {
	scene  scene.Scene
	camera *components.Camera
}

func (a *app) OnStart() {
	a.scene = scene.Load("main")
	a.scene.MakeCurrent()
	a.camera = a.scene.Root.Child("camera").Component("Camera").(*components.Camera)

	for k, v := range scene.MountMap {
		model.Load(v).Mount(k)
	}
}
func (a *app) OnUpdate() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())

	a.scene.Update()
	a.scene.Render()
}
func (a *app) OnExit() {}

func main() {
	window.Settings.SetVSync(true)
	window.Settings.SetTitle("MyGame")
	window.Settings.SetSize(1024, 768)

	engine.Start(&app{})
}
