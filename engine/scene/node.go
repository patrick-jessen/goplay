package scene

import (
	"encoding/json"
	"reflect"

	mgl "github.com/go-gl/mathgl/mgl32"
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
		name:           "root",
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

	c.Initialize(n)
	n.components[compName] = c
}

// Child returns the child with the given name.
// Returns nil if child does not exist.
func (n *Node) Child(name string) *Node {
	return n.children[name]
}

// Component returns the component with the given type.
// Returns nil if component does not exist.
func (n *Node) Component(name string) Component {
	return n.components[name]
}

func (n *Node) WorldTransform() mgl.Mat4 {
	return n.worldTransform
}

// update is called once every game loop.
func (n *Node) update() {
	if n.parent != nil {
		n.worldTransform = n.parent.worldTransform.Mul4(n.Transform.mat)
	}

	for _, c := range n.components {
		c.Update()
	}
	for _, c := range n.children {
		c.update()
	}
}

func (n *Node) render() {
	for _, c := range n.components {
		c.Render()
	}
	for _, c := range n.children {
		c.render()
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

// MarshalJSON encodes a node as JSON.
func (n *Node) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Transform  Transform            `json:"transform"`
		Children   map[string]*Node     `json:"children"`
		Components map[string]Component `json:"components"`
	}{
		Transform:  n.Transform,
		Children:   n.children,
		Components: n.components,
	}

	return json.Marshal(&tmp)
}

// String returns a JSON encoded string.
func (n Node) String() string {
	b, _ := json.Marshal(&n)
	return string(b)
}
