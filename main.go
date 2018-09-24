package main

import (
	"github.com/patrick-jessen/goplay/engine"
	"github.com/patrick-jessen/goplay/engine/texture"
	"github.com/patrick-jessen/goplay/engine/window"
)

func main() {
	window.Settings.SetVSync(true)
	window.Settings.SetTitle("MyGame")
	window.Settings.SetSize(1024, 768)

	texture.Settings.SetResolution(10)
	texture.Settings.Apply()

	engine.Start()
}
