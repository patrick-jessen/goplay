package main

import (
	"fmt"

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
	a.camera = a.scene.Root.Child("camera").Component("Camera").(*components.Camera)

	for k, v := range scene.MountMap {
		model.Load(v).Mount(k)
	}

	fmt.Println(a.scene.Root.Child("duck").Child("0"))
}
func (a *app) OnUpdate() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())

	a.scene.Update()
	a.scene.Render()
}
func (a *app) OnExit() {}

func main() {
	window.SetTitle("MyGame")
	window.SetVideoMode(false, 1024, 768)
	window.SetVerticalSync(true)
	engine.Start(&app{})
}
