package renderer

import "github.com/patrick-jessen/goplay/engine/scene"

type forwardRenderer struct {
}

func (f forwardRenderer) Render(s scene.Scene) {
	// Render shadow maps

	// Render scene normally

	// Render postprocessing quad
}
