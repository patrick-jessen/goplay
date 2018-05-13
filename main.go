package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/patrick-jessen/goplay/components"
	"github.com/patrick-jessen/goplay/engine"
	"github.com/patrick-jessen/goplay/engine/model"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/shader"
	"github.com/patrick-jessen/goplay/engine/texture"
	"github.com/patrick-jessen/goplay/engine/window"
)

type app struct {
	scene  scene.Scene
	camera *components.Camera
}

func (a *app) OnStart() {
	shader.Load("basic")
	texture.Load("default").Bind(0)

	a.scene = scene.New()
	child0 := a.scene.Root.NewChild("cubeNode")
	child1 := a.scene.Root.NewChild("cubeNode2")

	child2 := a.scene.Root.NewChild("camera")

	a.camera = &components.Camera{
		ProjectionMatrix: mgl.Perspective(mgl.DegToRad(45.0), float32(800)/float32(600), 0.1, 100.0),
	}

	child1.SetPosition(mgl.Vec3{1, 1, 0})
	child2.AddComponent(a.camera)
	child2.AddComponent(&components.ArcBall{
		Dist: 10,
	})

	model.Load("cube").Mount(child0)
	model.Load("cube").Mount(child1)

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
	engine.Start(&app{})
}
