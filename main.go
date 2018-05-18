package main

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/framebuffer"
	"github.com/patrick-jessen/goplay/engine/model"
	"github.com/patrick-jessen/goplay/engine/renderer"
	"github.com/patrick-jessen/goplay/engine/worker"

	"github.com/patrick-jessen/goplay/components"
	"github.com/patrick-jessen/goplay/engine"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/shader"
	"github.com/patrick-jessen/goplay/engine/window"
)

type quadMat struct {
	Shader shader.Shader
}

func (m quadMat) Apply() {
	m.Shader.Use()
}

var fb uint32

const msaa = 4

type app struct {
	scene     scene.Scene
	camera    *components.Camera
	postScene scene.Scene

	fb *framebuffer.FrameBuffer

	msaa *framebuffer.FrameBuffer

	renderer renderer.Renderer
}

func (a *app) OnStart() {
	gl.ClearColor(1, 1, 1, 1)
	a.postScene = scene.New()
	model.Load("quad").Mount(a.postScene.Root)

	s := shader.Load("fxaa")
	uRes := s.GetUniform("resolution")
	w, h := window.Settings.Size()
	gl.Uniform2f(uRes, float32(w), float32(h))

	a.postScene.Root.Child("0").Component("MeshRenderer").(*model.MeshRenderer).Mat = &quadMat{
		Shader: s,
	}

	a.scene = scene.Load("main")
	a.scene.MakeCurrent()
	a.camera = a.scene.Root.Child("camera").Component("Camera").(*components.Camera)

	for k, v := range scene.MountMap {
		model.Load(v).Mount(k)
	}

	a.fb = framebuffer.New(w, h, 1, 0)
	a.msaa = framebuffer.New(w, h, 1, 4)

	go func() {
		for {
			for i := 0; i < 3; i++ {
				locali := i
				worker.CallSynchronized(func() {
					aa = locali
					fmt.Println(locali)
				})
				<-time.After(2 * time.Second)
			}
		}

	}()

	a.renderer = renderer.NewForward()
	a.renderer.Initialize(a.scene)
}

var aa int

func (a *app) OnUpdate() {

	a.scene.Update() // <- move to engine
	shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())

	a.renderer.Render()
	return

	// MSAA
	if aa == 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())
		a.scene.Render()
	} else if aa == 1 {
		a.msaa.Bind()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())
		a.scene.Render()

		a.msaa.Blit(nil, 1024, 768, false)
	} else {
		a.fb.Bind()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())
		a.scene.Render()

		// Render post-processing quad
		framebuffer.Unbind()
		a.fb.BindColorTexture(0, 0)
		a.fb.BindDepthTexture(1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		a.postScene.Render()
	}
}
func (a *app) OnExit() {}

func main() {
	window.Settings.SetVSync(true)
	window.Settings.SetTitle("MyGame")
	window.Settings.SetSize(1024, 768)

	engine.Start(&app{})
}
