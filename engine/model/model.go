package model

import (
	"os"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/patrick-jessen/goplay/engine/model/geometry"
	"github.com/patrick-jessen/goplay/engine/model/gltf"
	"github.com/patrick-jessen/goplay/engine/scene"
	"github.com/patrick-jessen/goplay/engine/shader"
)

const modelDir = "./assets/models/"

var cache = make(map[string]Model)

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
			geom := gltf.GeometryFromPrimitive(m.file, &p)
			mr.geoms = append(mr.geoms, geom)
		}

		sn.AddComponent(mr)
	}

	for _, nidx := range gn.Children {
		gn := g.Nodes[nidx]
		child := sn.NewChild(gn.Name)
		m.mountChild(child, &gn)
	}
}

type MeshRenderer struct {
	node  *scene.Node
	geoms []*geometry.Geometry
}

func (mr *MeshRenderer) Initialize(n *scene.Node) {
	mr.node = n
}
func (mr *MeshRenderer) Update() {
}
func (mr *MeshRenderer) Render() {
	shader.SetModelMatrix(mr.node.WorldTransform())

	for _, g := range mr.geoms {
		g.Draw()
	}
}

// func loadNode(g *gltf.File, n *gltf.Node) *node.Node {
// 	node := node.New()
// 	node.Transform = mgl.Mat4{
// 		n.Matrix[0], n.Matrix[1], n.Matrix[2], n.Matrix[3],
// 		n.Matrix[4], n.Matrix[5], n.Matrix[6], n.Matrix[7],
// 		n.Matrix[8], n.Matrix[9], n.Matrix[10], n.Matrix[11],
// 		n.Matrix[12], n.Matrix[13], n.Matrix[14], n.Matrix[15],
// 	}
// 	t := mgl.Translate3D(n.Translation[0], n.Translation[1], n.Translation[2])
// 	r := mgl.QuatRotate(n.Rotation[3], mgl.Vec3{
// 		n.Rotation[0], n.Rotation[1], n.Rotation[2],
// 	}).Mat4()
// 	s := mgl.Scale3D(n.Scale[0], n.Scale[1], n.Scale[2])

// 	node.Transform = node.Transform.Mul4(t)
// 	node.Transform = node.Transform.Mul4(r)
// 	node.Transform = node.Transform.Mul4(s)

// 	// if n.Mesh >= 0 {

// 	// 	mr := &components.MeshRenderer{}

// 	// 	mesh := g.GlTF.Meshes[n.Mesh]
// 	// 	for _, p := range mesh.Primitives {
// 	// 		// Geometry
// 	// 		geom := gltf.GeometryFromPrimitive(g, &p)
// 	// 		mr.Geoms = append(mr.Geoms, geom)
// 	// 	}

// 	// 	node.AddComponent(mr)
// 	// }

// 	// for i, c := range n.Children {
// 	// 	child := loadNode(g, &g.GlTF.Nodes[c])
// 	// 	node.AddChild(fmt.Sprintf("%v", i), child)
// 	// }

// 	return node
// }
