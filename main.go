package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/shader"
	"github.com/patrick-jessen/goplay/engine/window"
)

type app struct {
	scene  scene.Scene
	shader shader.Shader
}

func (a *app) OnStart() {
	a.scene = scene.New()
	a.shader = shader.Load("basic")
}
func (a *app) OnUpdate() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	a.scene.Update()
	a.scene.Render()
}
func (a *app) OnExit() {}

func main() {
	window.SetTitle("MyGame")
	window.SetVideoMode(false, 1024, 768)
	engine.Start(&app{})
}
