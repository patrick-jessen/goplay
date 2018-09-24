// Package engine provides an interface for initializing, managing
// and deinitializing applications.
package engine

import (
	// Include components
	_ "github.com/patrick-jessen/goplay/components"
	"github.com/patrick-jessen/goplay/editor"
	"github.com/patrick-jessen/goplay/engine/renderer"
	"github.com/patrick-jessen/goplay/engine/resource"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/window"
	"github.com/patrick-jessen/goplay/engine/worker"
)

// Start starts the engine using the given application.
func Start() {
	go editor.Start()

	window.Create()
	defer window.Destroy()

	renderer.Initialize()
	defer renderer.Deinitialize()

	resource.LoadScene("main").MakeCurrent()

	for !window.ShouldClose() {
		window.Update()
		scene.Current().Update()
		renderer.Render()

		select {
		case ef := <-editor.Channel:
			ef()
		case work := <-worker.Channel:
			work()
		default:
		}
	}
}
