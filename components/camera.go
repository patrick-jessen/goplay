package components

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/patrick-jessen/goplay/engine/scene"
)

func init() {
	scene.RegisterComponent(&Camera{})
	scene.RegisterComponent(&ArcBall{})
}

type Camera struct {
	ProjectionMatrix mgl.Mat4
	node             *scene.Node
}

func New() *Camera {
	c := &Camera{
		ProjectionMatrix: mgl.Perspective(mgl.DegToRad(45.0), float32(800)/float32(600), 0.1, 100.0),
	}
	return c
}

func (c *Camera) Initialize(n *scene.Node) {
	c.node = n
}

func (c *Camera) Render() {}
func (c *Camera) Update() {}

func (c *Camera) ViewProjectionMatrix() mgl.Mat4 {
	return c.ProjectionMatrix.Mul4(c.node.WorldTransform())
}
