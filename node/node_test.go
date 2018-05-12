package node

import (
	"encoding/json"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
)

func init() {
	RegisterComponent(&testComponent{})
}

type testComponent struct {
	initializeCalled int
	updateCalled     int
	node             *Node
	Value            int `json:"value"`
}

func (t *testComponent) initialize(n *Node) {
	t.initializeCalled++
	t.node = n
}
func (t *testComponent) update() {
	t.updateCalled++
}

func TestNewNode(t *testing.T) {
	n := newNode()
	if n.mat != mgl.Ident4() {
		t.Errorf("wrong transformation:\n got %v\nexpected %v", n.mat, mgl.Ident4())
	}
	if n.worldTransform != mgl.Ident4() {
		t.Errorf("wrong world-transformation:\n got %v\nexpected %v", n.worldTransform, mgl.Ident4())
	}
	if n.children == nil {
		t.Error("children is nil-map")
	}
	if n.components == nil {
		t.Error("components is nil-map")
	}
}

func TestNodeNewChild(t *testing.T) {
	parent := newNode()
	returned := parent.NewChild("test")

	child, ok := parent.children["test"]
	if !ok {
		t.Fatal("child is not in child map")
	}
	if child.parent != parent {
		t.Error("child does not have the right parent")
	}
	if child.name != "test" {
		t.Errorf("incorrect name: got %v, expected %v", child.name, "test")
	}
	if returned != child {
		t.Error("child is not returned from NewChild()")
	}

	defer func() {
		recover()
	}()
	parent.NewChild("test")
	t.Error("child with same name can be added multiple times")
}

func TestNodeAddComponent(t *testing.T) {
	node := newNode()

	node.AddComponent(&testComponent{})

	comp, ok := node.components["testComponent"]
	if !ok {
		t.Error("component is not in component map")
	}
	tcomp, ok := comp.(*testComponent)
	if !ok {
		t.Error("component is not of the right type")
	}
	if tcomp.node != node {
		t.Error("node not set in component")
	}
	if tcomp.initializeCalled != 1 {
		t.Errorf("initialize() was not called the right number of times. got %v expected %v",
			tcomp.initializeCalled, 1)
	}

	defer func() {
		recover()
	}()
	node.AddComponent(&testComponent{})
	t.Error("component with same type can be added multiple times")
}

func TestNodeUpdate(t *testing.T) {
	parent := newNode()
	child := parent.NewChild("child")
	child.SetPosition(mgl.Vec3{1, 2, 3})
	child.AddComponent(&testComponent{})

	parent.update()

	if parent.worldTransform != mgl.Ident4() {
		t.Errorf("root has wrong world transform.\ngot %v\nexpected %v",
			parent.worldTransform, mgl.Ident4())
	}
	if child.worldTransform != mgl.Translate3D(1, 2, 3) {
		t.Errorf("child has wrong world transform.\ngot %v\nexpected %v",
			child.worldTransform, mgl.Translate3D(1, 2, 3))
	}
	comp := child.components["testComponent"].(*testComponent)
	if comp.updateCalled != 1 {
		t.Errorf("update() was not called the right number of times. got %v expected %v",
			comp.updateCalled, 1)
	}
}

func TestNodeJSONUnmarshal(t *testing.T) {
	jsonSrc := `{
		"transform": {
			"position":	[1,2,3],
			"rotation":	[4,5,6,7],
			"scale":		[8,9,10]
		},
		"components": {
			"testComponent": {
				"value": 42
			}
		},
		"children": {
			"test": {
				"transform": {
					"position":	[11,12,13]
				}
			}
		}
	}`

	n := newNode()
	e := json.Unmarshal([]byte(jsonSrc), n)
	if e != nil {
		t.Fatal(e)
	}

	if n.position != (mgl.Vec3{1, 2, 3}) ||
		n.rotation != (mgl.Quat{W: 4, V: mgl.Vec3{5, 6, 7}}) ||
		n.scale != (mgl.Vec3{8, 9, 10}) {
		t.Error("incorrect transform")
	}
	comp, ok := n.components["testComponent"]
	if !ok {
		t.Fatal("component not added to node")
	}
	tcomp, ok := comp.(*testComponent)
	if !ok {
		t.Fatal("component is not the right type")
	}
	if tcomp.Value != 42 {
		t.Errorf("component property not set. got %v, expected %v", tcomp.Value, 42)
	}
	child, ok := n.children["test"]
	if !ok {
		t.Fatal("child is not added to node")
	}
	if child.Transform.position != (mgl.Vec3{11, 12, 13}) {
		t.Error("child is not unmarshalled")
	}
}
