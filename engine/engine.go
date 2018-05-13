// Package engine provides an interface for initializing, managing
// and deinitializing applications.
package engine

import (
	"github.com/patrick-jessen/goplay/engine/window"
)

// Application is the base of the application.
type Application interface {
	OnStart()
	OnUpdate()
	OnExit()
}

// Start starts the engine using the given application.
func Start(a Application) {
	defer window.Deinitialize()

	window.Create()
	defer window.Close()

	a.OnStart()
	defer a.OnExit()

	for !window.ShouldClose() {
		window.Update()
		a.OnUpdate()
	}
}
