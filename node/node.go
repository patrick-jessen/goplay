package node

import (
	"encoding/json"
	"reflect"

	mgl "github.com/go-gl/mathgl/mgl64"
)

// Node is a node in the 3D scene graph.
type Node struct {
	Transform
	children   map[string]*Node
	components map[string]Component

	parent         *Node
	name           string
	worldTransform mgl.Mat4
}

// newNode creates a new node.
func newNode() *Node {
	return &Node{
		Transform:      newTransform(),
		children:       make(map[string]*Node),
		components:     make(map[string]Component),
		worldTransform: mgl.Ident4(),
	}
}

// initialize is called whenever a node is attached to a (new) parent.
func (n *Node) initialize(parent *Node, name string) {
	n.parent = parent
	n.name = name
}

// NewChild creates a new child on the current node.
// The child is returned.
func (n *Node) NewChild(name string) *Node {
	if _, ok := n.children[name]; ok {
		panic("child already exists: " + name)
	}

	child := newNode()
	child.initialize(n, name)
	n.children[name] = child
	return child
}

// AddComponent adds a component.
// A node can only have one instance of the same component.
func (n *Node) AddComponent(c Component) {
	compName := reflect.TypeOf(c).Elem().Name()
	if _, ok := n.components[compName]; ok {
		panic("component already exists: " + compName)
	}

	c.initialize(n)
	n.components[compName] = c
}

// update is called once every game loop.
func (n *Node) update() {
	if n.parent != nil {
		n.worldTransform = n.parent.worldTransform.Mul4(n.Transform.mat)
	}

	for _, c := range n.components {
		c.update()
	}
	for _, c := range n.children {
		c.update()
	}
}

// UnmarshalJSON decodes a node from JSON.
func (n *Node) UnmarshalJSON(data []byte) error {
	var objMap map[string]*json.RawMessage
	e := json.Unmarshal(data, &objMap)
	if e != nil {
		return e
	}

	if t, ok := objMap["transform"]; ok {
		json.Unmarshal(*t, &n.Transform)
	}
	if c, ok := objMap["children"]; ok {
		var childMap map[string]*json.RawMessage
		e = json.Unmarshal(*c, &childMap)
		if e != nil {
			return e
		}

		for k, v := range childMap {
			var child Node
			e = json.Unmarshal(*v, &child)
			if e != nil {
				return e
			}
			n.children[k] = &child
		}
	}
	if c, ok := objMap["components"]; ok {
		var compMap map[string]*json.RawMessage
		e = json.Unmarshal(*c, &compMap)
		if e != nil {
			return e
		}

		for k, v := range compMap {
			typ, ok := componentMap[k]
			if !ok {
				panic("invalid component type: " + k)
			}
			comp := reflect.New(typ)
			e = json.Unmarshal(*v, comp.Interface())
			if e != nil {
				return e
			}
			n.components[k] = comp.Interface().(Component)
		}
	}

	return nil
}
