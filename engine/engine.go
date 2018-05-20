// Package engine provides an interface for initializing, managing
// and deinitializing applications.
package engine

import (
	"github.com/patrick-jessen/goplay/editor"
	"github.com/patrick-jessen/goplay/engine/renderer"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/window"
	"github.com/patrick-jessen/goplay/engine/worker"
)

// Application is the base of the application.
type Application interface {
	OnStart()
	OnUpdate()
	OnExit()
}

// Start starts the engine using the given application.
func Start(a Application) {
	go editor.Start()

	window.Create()
	defer window.Destroy()

	renderer.Initialize()
	defer renderer.Deinitialize()

	a.OnStart()
	defer a.OnExit()

	for !window.ShouldClose() {
		window.Update()
		a.OnUpdate()

		scene.Current().Update()
		renderer.Render()

		select {
		case ef := <-editor.EditorChannel:
			ef()
		case wf := <-worker.Channel:
			wf()
		default:
		}
	}
}
