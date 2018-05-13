package scene

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

func (t *testComponent) Initialize(n *Node) {
	t.initializeCalled++
	t.node = n
}
func (t *testComponent) Update() {
	t.updateCalled++
}

func Test_newNode(t *testing.T) {
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
	if n.name != "root" {
		t.Error("root name not set")
	}
}

func TestNode_initialize(t *testing.T) {
	parent := newNode()
	child := newNode()

	child.initialize(parent, "child")

	if child.parent != parent {
		t.Errorf("incorrect parent. got %v, expected %v", child.parent, parent)
	}
	if child.name != "child" {
		t.Errorf("incorrect name. got %v, expected %v", child.name, "child")
	}
}

func TestNode_NewChild(t *testing.T) {
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

func TestNode_AddComponent(t *testing.T) {
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

func TestNode_Child(t *testing.T) {
	n := newNode()
	c := n.NewChild("child")

	if n.Child("child") != c {
		t.Errorf("not the right child. got %v, expected %v", n.Child("child"), c)
	}
}

func TestNode_Component(t *testing.T) {
	n := newNode()
	n.AddComponent(&testComponent{})

	comp := n.Component("testComponent")
	_, ok := comp.(*testComponent)
	if !ok {
		t.Errorf("not the right component. got %v", comp)
	}
}

func TestNode_update(t *testing.T) {
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

func TestNode_UnmarshalJSON(t *testing.T) {
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

	// Test for non-existing component
	jsonSrc = `{
		"components": {
			"doesNotExist": {
				"value": 42
			}
		}
	}`
	n = newNode()
	defer func() {
		recover()
	}()

	e = json.Unmarshal([]byte(jsonSrc), n)
	if e != nil {
		t.Fatal(e)
	}
	t.Error("non-existing component types allowed")
}

func TestNode_MarshalJSON(t *testing.T) {
	expected := `{"transform":{"position":[0,0,0],"rotation":[1,0,0,0],"scale":[1,1,1]},"children":{"child":{"transform":{"position":[0,0,0],"rotation":[1,0,0,0],"scale":[1,1,1]},"children":{},"components":{"testComponent":{"value":1234}}}},"components":{}}`

	n := newNode()
	n.NewChild("child").AddComponent(&testComponent{Value: 1234})

	b, e := json.Marshal(n)
	if e != nil {
		t.Fatal("failed to marshal node")
	}
	if string(b) != expected {
		t.Errorf("did not marshal correctly.\n got %v\n expected %v", string(b), expected)
	}
}

func TestNode_String(t *testing.T) {
	expected := `{"transform":{"position":[0,0,0],"rotation":[1,0,0,0],"scale":[1,1,1]},"children":{"child":{"transform":{"position":[0,0,0],"rotation":[1,0,0,0],"scale":[1,1,1]},"children":{},"components":{"testComponent":{"value":1234}}}},"components":{}}`

	n := newNode()
	n.NewChild("child").AddComponent(&testComponent{Value: 1234})

	if n.String() != expected {
		t.Errorf("did not marshal correctly.\n got %v\n expected %v", n.String(), expected)
	}
}
