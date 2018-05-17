package components

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/window"
)

func init() {
	scene.RegisterComponent(&Camera{})
	scene.RegisterComponent(&ArcBall{})
}

type Camera struct {
	FOV              float32
	ProjectionMatrix mgl.Mat4
	node             *scene.Node
}

func New() *Camera {
	c := &Camera{
		FOV: 45,
	}
	return c
}

func (c *Camera) Initialize(n *scene.Node) {
	c.node = n

	w, h := window.Settings.Size()
	c.ProjectionMatrix = mgl.Perspective(
		mgl.DegToRad(c.FOV),
		float32(w)/float32(h),
		0.01, 1000.0)

	window.AddResizeHandler(func(w, h int) {
		c.ProjectionMatrix = mgl.Perspective(
			mgl.DegToRad(c.FOV),
			float32(w)/float32(h),
			0.01, 1000.0)
	})
}

func (c *Camera) Render() {}
func (c *Camera) Update() {}

func (c *Camera) ViewProjectionMatrix() mgl.Mat4 {
	return c.ProjectionMatrix.Mul4(c.node.WorldTransform())
}
