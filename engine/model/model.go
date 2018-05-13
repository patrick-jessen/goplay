package model

import (
	"fmt"
	"os"

	"github.com/patrick-jessen/goplay/engine/texture"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/patrick-jessen/goplay/engine/material"
	"github.com/patrick-jessen/goplay/engine/model/geometry"
	"github.com/patrick-jessen/goplay/engine/model/gltf"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/shader"
)

const modelDir = "./assets/models/"

var cache = make(map[string]Model)
var nameIter = 0

// Load returns a model by either loading it or reading from cache.
func Load(name string) Model {
	// Read form cache
	if val, ok := cache[name]; ok {
		return val
	}
	// Load from disk
	cache[name] = loadModel(name)
	return cache[name]
}

// loadModel loads a model from a file.
func loadModel(name string) Model {
	file := modelDir + name
	if _, err := os.Stat(file + ".glb"); err == nil {
		file += ".glb"
	} else if _, err := os.Stat(file + ".gltf"); err == nil {
		file += ".gltf"
	} else {
		panic("Model not found: " + name)
	}

	return Model{file: gltf.Load(file)}
}

type Model struct {
	file *gltf.File
}

func (m Model) Mount(sn *scene.Node) {
	g := m.file.GlTF
	scene := g.Scenes[g.Scene]

	pos := sn.Position()
	rot := sn.Rotation()
	scal := sn.Scale()

	nameIter = 0
	m.mountChild(sn, &gltf.Node{
		Translation: []float32{pos.X(), pos.Y(), pos.Z()},
		Rotation:    []float32{rot.W, rot.V.X(), rot.V.Y(), rot.V.Z()},
		Scale:       []float32{scal.X(), scal.Y(), scal.Z()},
		Children:    scene.Nodes,
		Mesh:        -1,
	})
}

func (m Model) mountChild(sn *scene.Node, gn *gltf.Node) {
	g := m.file.GlTF

	sn.SetPosition(mgl.Vec3{gn.Translation[0], gn.Translation[1], gn.Translation[2]})
	sn.SetRotation(mgl.Quat{W: gn.Rotation[0], V: mgl.Vec3{gn.Rotation[1], gn.Rotation[2], gn.Rotation[3]}})
	sn.SetScale(mgl.Vec3{gn.Scale[0], gn.Scale[1], gn.Scale[2]})

	if gn.Mesh >= 0 {
		mr := &MeshRenderer{}
		mesh := g.Meshes[gn.Mesh]
		for _, p := range mesh.Primitives {
			// Set geometry
			geom := gltf.GeometryFromPrimitive(m.file, &p)
			mr.geoms = append(mr.geoms, geom)

			// Set material
			gmat := g.Materials[p.Material]
			mat := material.New()
			if gmat.PbrMetallicRoughness.BaseColorTexture.Index >= 0 {
				t := g.Textures[gmat.PbrMetallicRoughness.BaseColorTexture.Index]
				tsrc := g.Images[t.Source]
				mat.Textures[0] = texture.Load(tsrc.URI)
			}
			mr.mats = append(mr.mats, mat)
		}

		sn.AddComponent(mr)
	}

	for _, nidx := range gn.Children {
		gn := g.Nodes[nidx]

		nam := gn.Name
		if len(nam) == 0 || sn.NewChild(gn.Name) != nil {
			nam = fmt.Sprintf("%v", nameIter)
			nameIter++
		}
		child := sn.NewChild(nam)
		m.mountChild(child, &gn)
	}
}

type MeshRenderer struct {
	node  *scene.Node
	geoms []*geometry.Geometry
	mats  []material.Material
}

func (mr *MeshRenderer) Initialize(n *scene.Node) {
	mr.node = n
}
func (mr *MeshRenderer) Update() {
}
func (mr *MeshRenderer) Render() {
	shader.SetModelMatrix(mr.node.WorldTransform())

	for i, g := range mr.geoms {
		mr.mats[i].Apply()
		g.Draw()
	}
}
