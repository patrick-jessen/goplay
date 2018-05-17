// Package engine provides an interface for initializing, managing
// and deinitializing applications.
package engine

import (
	"github.com/patrick-jessen/goplay/editor"
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

	a.OnStart()
	defer a.OnExit()

	for !window.ShouldClose() {
		window.Update()
		a.OnUpdate()

		select {
		case ef := <-editor.EditorChannel:
			ef()
		case wf := <-worker.Channel:
			wf()
		default:
		}
	}
}
