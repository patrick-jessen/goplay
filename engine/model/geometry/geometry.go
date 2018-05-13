package geometry

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

// Geometry represents renderable geometry.
type Geometry struct {
	handle     uint32 // Handle to OpenGL VertexArray.
	numIndices int32  // The number of indices.

	PrimType       uint32 // The type of primitives to render.
	IndexBuffer    Buffer
	PositionBuffer Buffer
	TexCoordBuffer Buffer
	NormalBuffer   Buffer
	TangentBuffer  Buffer
	hasIndices     bool
}

// componentSizeFromType gets the size of the component type.
// e.g. UNSIGNED_INT is 4 bytes.
func componentSizeFromType(compType uint32) int32 {
	switch compType {
	case gl.BYTE:
		fallthrough
	case gl.UNSIGNED_BYTE:
		return 1
	case gl.SHORT:
		fallthrough
	case gl.UNSIGNED_SHORT:
		return 2
	case gl.UNSIGNED_INT:
		return 4
	case gl.FLOAT:
		return 4
	default:
		panic("Unexpected componentType: " + string(compType))
	}
}

// Initialize initializes the geometry and uploads data to the GPU.
func (g *Geometry) Initialize() {
	if g.handle != 0 {
		return
	}

	// Determine if geometry has an index buffer
	if g.IndexBuffer.ComponentType != 0 {
		g.hasIndices = true
	}

	// Calculate number of indices
	if g.hasIndices {
		compSize := componentSizeFromType(g.IndexBuffer.ComponentType)
		g.numIndices = int32(len(g.IndexBuffer.Data)) / compSize
	} else {
		compSize := componentSizeFromType(g.PositionBuffer.ComponentType)
		g.numIndices = int32(len(g.PositionBuffer.Data)) / compSize
	}

	// Set buffer targets
	g.IndexBuffer.target = gl.ELEMENT_ARRAY_BUFFER
	g.PositionBuffer.target = gl.ARRAY_BUFFER
	g.NormalBuffer.target = gl.ARRAY_BUFFER
	g.TexCoordBuffer.target = gl.ARRAY_BUFFER
	g.TangentBuffer.target = gl.ARRAY_BUFFER

	// Initialize buffers
	g.IndexBuffer.initialize()
	g.PositionBuffer.initialize()
	g.NormalBuffer.initialize()
	g.TexCoordBuffer.initialize()
	g.TangentBuffer.initialize()

	// Create and bind VertexArray
	gl.GenVertexArrays(1, &g.handle)
	gl.BindVertexArray(g.handle)

	// Bind/enable buffers within the VertexArray
	g.IndexBuffer.bind()
	g.PositionBuffer.enable(0)
	g.NormalBuffer.enable(1)
	g.TexCoordBuffer.enable(2)
	g.TangentBuffer.enable(3)

	gl.BindVertexArray(0)
}

// Free frees the resources of the geometry
func (g *Geometry) Free() {
	gl.DeleteVertexArrays(1, &g.handle)
	g.IndexBuffer.free()
	g.PositionBuffer.free()
	g.NormalBuffer.free()
	g.TexCoordBuffer.free()
}

// Draw draws the geometry
func (g *Geometry) Draw() {
	gl.BindVertexArray(g.handle)
	if g.hasIndices {
		gl.DrawElements(
			g.PrimType,
			g.numIndices,
			g.IndexBuffer.ComponentType,
			gl.PtrOffset(g.IndexBuffer.ByteOffset),
		)
	} else {
		gl.DrawArrays(g.PrimType, 0, g.numIndices)
	}
}
