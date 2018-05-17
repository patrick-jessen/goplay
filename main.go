package main

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/framebuffer"
	"github.com/patrick-jessen/goplay/engine/log"
	"github.com/patrick-jessen/goplay/engine/model"
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

func initMSAA() {
	w, h := window.Settings.Size()
	log.Info("Msaa", "w", w, "h", h)

	gl.GenFramebuffers(1, &fb)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb)

	var rb uint32
	gl.GenRenderbuffers(1, &rb)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rb)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa, gl.RGB8, int32(w), int32(h))
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	var rbd uint32
	gl.GenRenderbuffers(1, &rbd)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rbd)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, msaa, gl.DEPTH_COMPONENT, int32(w), int32(h))
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, rb)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, rbd)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		log.Info("", "res", gl.CheckFramebufferStatus(gl.FRAMEBUFFER))
		panic("ERROR::FRAMEBUFFER:: Framebuffer is not complete!")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func preRenderMSAA() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb)
}

func postRenderMSAA() {
	w, h := window.Settings.Size()

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fb)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(0, 0, int32(w), int32(h), 0, 0, int32(w), int32(h), gl.COLOR_BUFFER_BIT, gl.NEAREST)
}

type app struct {
	scene     scene.Scene
	camera    *components.Camera
	postScene scene.Scene

	fb *framebuffer.FrameBuffer
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

	a.fb = framebuffer.New(1)

	initMSAA()

	go func() {
		for {
			for i := 0; i < 3; i++ {
				<-time.After(2 * time.Second)
				locali := i
				worker.CallSynchronized(func() {
					aa = locali
					fmt.Println(locali)
				})
			}
		}

	}()
}

var aa int

func (a *app) OnUpdate() {

	a.scene.Update() // <- move to engine

	// Render shadow maps
	// TODO

	// Render normally

	////sssssssssssssssssssssss
	// framebuffer.Use(a.fb)
	// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())
	// a.scene.Render()

	// // Render post-processing quad

	// framebuffer.Use(nil)
	// a.fb.BindColorTexture(0, 0)
	// a.fb.BindDepthTexture(1)
	// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// a.postScene.Render()
	//s ssssssssssssssssssssssssssss

	// MSAA
	if aa == 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())
		a.scene.Render()
	} else if aa == 1 {
		preRenderMSAA()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())
		a.scene.Render()

		postRenderMSAA()
	} else {
		framebuffer.Use(a.fb)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shader.SetViewProjectionMatrix(a.camera.ViewProjectionMatrix())
		a.scene.Render()

		// Render post-processing quad
		framebuffer.Use(nil)
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
