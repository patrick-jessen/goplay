package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine"
	"github.com/patrick-jessen/goplay/engine/window"
)

type app struct{}

func (a *app) OnStart() {
	gl.ClearColor(1, 0, 0, 1)
}
func (a *app) OnUpdate() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
func (a *app) OnExit() {}

func main() {
	window.SetTitle("MyGame")
	window.SetVideoMode(false, 1024, 768)
	engine.Start(&app{})
}
