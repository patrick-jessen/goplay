// Package shader implements OpenGL shader programs.
// Shaders should be located under assets/shaders/{name}/.
// Each shader folder should have a {name}.vert and a {name}.frag.
package shader

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl64"
)

const shaderDir = "./assets/shaders/"

var cache = make(map[string]Shader)
var ubo uint32

// Load returns a shader by either loading it or reading from cache.
func Load(name string) Shader {
	// Read form cache
	if val, ok := cache[name]; ok {
		return val
	}
	// Load from disk
	cache[name] = Shader{handle: loadProgram(name)}
	return cache[name]
}

// Shader represents an OpenGL shader program.
type Shader struct {
	handle uint32
}

// Use sets a shader program for use.
func (s Shader) Use() {
	gl.UseProgram(s.handle)
}

func (s Shader) GetUniform(name string) int32 {
	s.Use()
	return gl.GetUniformLocation(s.handle, gl.Str(name+"\x00"))
}

// SetViewProjectionMatrix sets the view-projection matrix for all shaders.
func SetViewProjectionMatrix(m mgl.Mat4) {
	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 64, gl.Ptr(&m[0]))
}

// SetModelMatrix sets the model matrix for all shaders.
func SetModelMatrix(m mgl.Mat4) {
	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 64, 64, gl.Ptr(&m[0]))
}

// loadProgram loads shaders from files and creates a shader program.
func loadProgram(name string) uint32 {
	file := shaderDir + name + "/" + name

	vertSrc, e := ioutil.ReadFile(file + ".vert")
	if e != nil {
		panic("failed to load vertex shader:\n" + e.Error())
	}
	fragSrc, e := ioutil.ReadFile(file + ".frag")
	if e != nil {
		panic("failed to load fragment shader:\n" + e.Error())
	}

	vert := compileShader(gl.VERTEX_SHADER, string(vertSrc))
	frag := compileShader(gl.FRAGMENT_SHADER, string(fragSrc))

	handle := gl.CreateProgram()

	gl.AttachShader(handle, vert)
	gl.AttachShader(handle, frag)

	linkProgram(handle)

	gl.DeleteShader(vert)
	gl.DeleteShader(frag)

	// TEMP
	initializeUniformBuffer()
	gl.UseProgram(handle)
	ubi := gl.GetUniformBlockIndex(handle, gl.Str("shader_data\x00"))
	gl.UniformBlockBinding(handle, ubi, 0)

	// Other uniforms
	gl.Uniform1i(gl.GetUniformLocation(handle, gl.Str("tex0\x00")), 0)

	return handle
}

func linkProgram(handle uint32) {
	gl.LinkProgram(handle)

	var status int32
	gl.GetProgramiv(handle, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(handle, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(handle, logLength, nil, gl.Str(log))

		panic("failed to link program:\n" + log)
	}
}

func compileShader(t uint32, src string) uint32 {
	// Make sure string is null-terminated
	src += "\x00"

	csrc, free := gl.Strs(src)
	defer free()

	handle := gl.CreateShader(t)
	gl.ShaderSource(handle, 1, csrc, nil)
	gl.CompileShader(handle)

	var status int32
	gl.GetShaderiv(handle, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(handle, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(handle, logLength, nil, gl.Str(log))

		panic("failed to compile shader:\n" + log)
	}
	return handle
}

func initializeUniformBuffer() {
	if ubo != 0 {
		return
	}

	proj := mgl.Perspective(mgl.DegToRad(45.0), float64(800)/float64(600), 0.1, 100.0)
	view := mgl.Translate3D(0, 0, -10)
	model := mgl.Ident4()

	buf := bytes.NewBuffer([]byte{})

	dat := struct {
		ViewProj mgl.Mat4
		Model    mgl.Mat4
	}{
		ViewProj: proj.Mul4(view),
		Model:    model,
	}
	binary.Write(buf, binary.LittleEndian, &dat)

	var handle uint32
	gl.GenBuffers(1, &handle)
	gl.BindBuffer(gl.UNIFORM_BUFFER, handle)
	gl.BufferData(gl.UNIFORM_BUFFER, 128, gl.Ptr(buf.Bytes()), gl.STATIC_DRAW)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, 0, handle)

	ubo = handle
}
