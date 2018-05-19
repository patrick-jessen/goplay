package components

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/patrick-jessen/goplay/engine/log"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/shader"
	"github.com/patrick-jessen/goplay/engine/window"
)

type ArcBall struct {
	RotX float32
	RotY float32
	Dist float32

	zoomVelocity float32
	rotXVelocity float32
	rotYVelocity float32

	node *scene.Node
}

const damping = 0.90

func (c *ArcBall) Initialize(n *scene.Node) {
	c.node = n

	log.Warn("ArcBall.Update currently calls shader.SetViewPosition(). Why can't this be done elsewhere?")
}
func (c *ArcBall) Render() {}

func (c *ArcBall) Update() {

	// Handle movement
	if window.MouseButton(0) {
		move := window.MouseMove().Mul(0.005)

		if move.X()*move.X() < c.rotYVelocity*c.rotYVelocity {
			c.rotYVelocity *= damping
		} else if move.X() != c.rotYVelocity {
			c.rotYVelocity += 0.1 * move.X()
		}

		if move.Y()*move.Y() < c.rotXVelocity*c.rotXVelocity {
			c.rotXVelocity *= damping
		} else if move.Y() != c.rotXVelocity {
			c.rotXVelocity += 0.1 * move.Y()
		}

	} else {
		c.rotXVelocity *= damping
		c.rotYVelocity *= damping
	}

	c.RotX += c.rotXVelocity
	c.RotY += c.rotYVelocity

	if c.RotX > math.Pi/2-0.001 {
		c.RotX = math.Pi/2 - 0.001
	} else if c.RotX < -math.Pi/2+0.001 {
		c.RotX = -math.Pi/2 + 0.001
	}

	// Handle zoom
	if window.MouseScroll() == 0 {
		c.zoomVelocity *= damping
	} else {
		c.zoomVelocity += window.MouseScroll() * c.Dist * 0.005
	}
	c.Dist -= c.zoomVelocity

	// Calculate camera position
	sinY := float32(math.Sin(float64(c.RotY)))
	cosY := float32(math.Cos(float64(c.RotY)))

	sinX := float32(math.Sin(float64(c.RotX)))
	cosX := float32(math.Cos(float64(c.RotX)))

	pos := mgl.Vec3{
		c.Dist * sinY * -cosX,
		c.Dist * sinX,
		c.Dist * -cosY * -cosX,
	}

	view := mgl.LookAtV(pos, mgl.Vec3{}, mgl.Vec3{0, 1, 0})
	c.node.SetMatrix(view)

	shader.SetViewPosition(pos)
}
