package geometry

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

// Buffer represents a data buffer on the GPU.
type Buffer struct {
	handle uint32 // The handle to the OpenGL buffer.
	target uint32 // The target that the GPU buffer should be bound to.

	ByteOffset    int    // The offset relative to the start of the data in bytes.
	ByteStride    int32  // The stride, in bytes, between vertex attributes.
	ComponentType uint32 // The datatype of components in the attribute.
	Normalized    bool   // Specifies whether integer data values should be normalized.
	NumComponents int32  // The number of components per vertex attribute.
	Data          []byte // The buffer data
}

// initialize uploads data to the GPU.
// After calling this function, b.Data will be nil.
func (b *Buffer) initialize() {
	if b.Data == nil {
		return
	}

	// Create buffer object
	gl.GenBuffers(1, &b.handle)

	// Bind buffer and upload data to it
	gl.BindBuffer(b.target, b.handle)
	gl.BufferData(b.target, len(b.Data), gl.Ptr(b.Data), gl.STATIC_DRAW)
	gl.BindBuffer(b.target, 0)

	b.Data = nil
}

// enable binds the buffer and enables vertex attribute.
// Intended to be used with VertexArray.
func (b *Buffer) enable(vertexAttribIndex uint32) {
	if b.handle == 0 {
		return
	}

	gl.BindBuffer(b.target, b.handle)
	gl.VertexAttribPointer(
		vertexAttribIndex,
		b.NumComponents,
		b.ComponentType,
		b.Normalized,
		b.ByteStride,
		gl.PtrOffset(b.ByteOffset),
	)
	gl.EnableVertexAttribArray(vertexAttribIndex)
}

// bind binds the buffer to its target.
func (b *Buffer) bind() {
	if b.handle == 0 {
		return
	}

	gl.BindBuffer(b.target, b.handle)
}

// free frees the buffer.
func (b *Buffer) free() {
	if b.handle == 0 {
		return
	}

	gl.DeleteBuffers(1, &b.handle)
}
