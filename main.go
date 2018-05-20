package main

import (
	"github.com/go-gl/gl/v3.2-core/gl"

	"github.com/patrick-jessen/goplay/engine"
	"github.com/patrick-jessen/goplay/engine/log"
	"github.com/patrick-jessen/goplay/engine/model"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/window"

	_ "github.com/patrick-jessen/goplay/components"
)

type app struct{}

func (a *app) OnStart() {
	gl.ClearColor(1, 1, 1, 1)

	scene.Load("main").MakeCurrent()
	log.Warn("app.OnStart() currently loads mounted models. This should be done elsewhere")
	for k, v := range scene.MountMap {
		model.Load(v).Mount(k)
	}
}

func (a *app) OnUpdate() {}
func (a *app) OnExit()   {}

func main() {
	window.Settings.SetVSync(true)
	window.Settings.SetTitle("MyGame")
	window.Settings.SetSize(1024, 768)

	engine.Start(&app{})
}
