package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine"
	"github.com/patrick-jessen/goplay/engine/model"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/shader"
	"github.com/patrick-jessen/goplay/engine/texture"
	"github.com/patrick-jessen/goplay/engine/window"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type app struct {
	scene   scene.Scene
	shader  shader.Shader
	texture texture.Texture
}

func (a *app) OnStart() {
	shader.Load("basic")
	a.texture = texture.Load("default")

	a.texture.Bind(0)

	a.scene = scene.New()
	cn := a.scene.Root.NewChild("cubeNode")
	cn2 := a.scene.Root.NewChild("cubeNode2")

	cn2.SetPosition(mgl.Vec3{1, 0, 0})

	model.Load("cube").Mount(cn)
	model.Load("cube").Mount(cn2)
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
