// Package gltf implements a glTF 2.0 loader.
// The package does NOT peform any validation.
package gltf

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/patrick-jessen/goplay/engine/model/geometry"
)

// File represents the binary glTF format
type File struct {
	GlTF     GlTF     // The glTF file
	Chunks   [][]byte // Only available when loading .glb
	Location string   // Directory of the file.
	File     string   // Path to the file.
}

// Load loads a glTF file.
// Accepts both .gltf and .glb.
func Load(file string) *File {
	data := readFile(file)
	out := &File{
		Location: filepath.Dir(file),
		File:     file,
	}

	if strings.HasSuffix(file, ".gltf") {
		json.Unmarshal(data, &out.GlTF)
	} else if strings.HasSuffix(file, ".glb") {
		iter := 0
		readGlbHeader(&iter, data)
		for iter < len(data) {
			d, t := readGlbChunk(&iter, data)
			switch t {
			case typeJSON:
				json.Unmarshal(d, &out.GlTF)
			case typeBIN:
				out.Chunks = append(out.Chunks, d)
			default:
				panic("Invalid chunk type")
			}
		}
	} else {
		panic("Invalid file format")
	}
	return out
}

const (
	magicValue   = 0x46546C67
	versionValue = 0x2
	typeJSON     = 0x4E4F534A
	typeBIN      = 0x004E4942
)

// readFile reads all bytes from a file.
func readFile(file string) []byte {
	data, e := ioutil.ReadFile(file)
	if e != nil {
		panic(e)
	}
	return data
}

// readInt reads a single 32 bit unsigned int from the file
func readInt(iter *int, d []byte) (res uint) {
	res = uint(binary.LittleEndian.Uint32(d[*iter:]))
	*iter += 4
	return
}

// readBytes reads a number of bytes from the file
func readBytes(iter *int, d []byte, len uint) (res []byte) {
	res = d[*iter : *iter+int(len)]
	*iter += int(len)
	return
}

// readGlbHeader reads and verifies the .glb header
// Header format: [Magic:u32, Version:u32, FileLen:u32]
func readGlbHeader(iter *int, d []byte) {
	// Make sure this is glTF format
	if readInt(iter, d) != magicValue {
		panic("File is not of glTF format")
	}
	// Make sure this is the right version of the format
	if readInt(iter, d) != versionValue {
		panic("glTF file is not the right version")
	}
	// Read length
	_ = readInt(iter, d)
}

// readGlbChunk reads a data chunk
// Chunk format: [dataLen:u32, type:u32, data:u8[]]
func readGlbChunk(iter *int, d []byte) ([]byte, uint) {
	length := readInt(iter, d)
	ctype := readInt(iter, d)
	data := readBytes(iter, d, length)
	return data, ctype
}

// bufferFromAccessor creates a buffer object from a gltf accessor.
func bufferFromAccessor(g *File, a *Accessor) geometry.Buffer {
	bufView := &g.GlTF.BufferViews[a.BufferView]
	data := dataFromBufferView(g, bufView)

	return geometry.Buffer{
		ByteOffset:    int(a.ByteOffset),
		ByteStride:    int32(bufView.ByteStride),
		ComponentType: uint32(a.ComponentType),
		Normalized:    a.Normalized,
		NumComponents: numComponentsInType(a.Type),
		Data:          data,
	}
}

// numComponentsInType returns the number of components in a type.
// e.g. VEC2 has 2 components.
func numComponentsInType(t string) int32 {
	switch t {
	case "SCALAR":
		return 1
	case "VEC2":
		return 2
	case "VEC3":
		return 3
	case "VEC4":
		return 4
	case "MAT2":
		return 4
	case "MAT3":
		return 9
	case "MAT4":
		return 16
	default:
		panic("Invalid data type: " + t)
	}
}

// dataFromBufferView returns the data associated with a buffer view.
func dataFromBufferView(g *File, b *BufferView) []byte {
	buffer := g.GlTF.Buffers[b.Buffer]
	var data []byte
	var e error

	if len(buffer.URI) == 0 {
		// Load from blob
		data = g.Chunks[b.Buffer]

	} else if strings.HasPrefix(buffer.URI, "data:") {
		// Load from data URI
		idx := strings.Index(buffer.URI, ";base64,") + len(";base64,")
		data, e = base64.StdEncoding.DecodeString(buffer.URI[idx:])
		if e != nil {
			panic(e)
		}

	} else {
		// Load from file
		data, e = ioutil.ReadFile(g.Location + "/" + buffer.URI)
		if e != nil {
			panic(e)
		}
	}

	return data[b.ByteOffset : b.ByteOffset+b.ByteLength]
}

// GeometryFromPrimitive creates a geometry object from a glTF primitive.
func GeometryFromPrimitive(g *File, prim *MeshPrimitive) *geometry.Geometry {
	// Create geometry
	geom := &geometry.Geometry{
		PrimType: uint32(prim.Mode),
	}

	// Specify index buffer
	if prim.Indices >= 0 {
		indexAccessor := &g.GlTF.Accessors[prim.Indices]
		geom.IndexBuffer = bufferFromAccessor(g, indexAccessor)
	}

	// Specify attribute buffers
	for key, val := range prim.Attributes {
		accessor := &g.GlTF.Accessors[val]

		strs := strings.Split(key, "_")
		switch strs[0] {
		case "POSITION":
			geom.PositionBuffer = bufferFromAccessor(g, accessor)
		case "NORMAL":
			geom.NormalBuffer = bufferFromAccessor(g, accessor)
		case "TANGENT":
			geom.TangentBuffer = bufferFromAccessor(g, accessor)
		case "TEXCOORD":
			if strs[1] == "0" {
				geom.TexCoordBuffer = bufferFromAccessor(g, accessor)
			}
		case "COLOR":
		case "JOINTS":
		case "WEIGHTS":
		default:
		}
	}

	geom.Initialize()
	return geom
}
