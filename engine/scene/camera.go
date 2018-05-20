package scene

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/patrick-jessen/goplay/engine/window"
)

func init() {
	RegisterComponent(&Camera{})
}

type Camera struct {
	FOV              float32
	ProjectionMatrix mgl.Mat4
	node             *Node
}

func NewCamera() *Camera {
	c := &Camera{
		FOV: 45,
	}
	return c
}

func (c *Camera) Initialize(n *Node) {
	c.node = n
	if n.scene.camera == nil {
		n.scene.camera = c
	}

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
