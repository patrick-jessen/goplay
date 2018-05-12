package node

import (
	"encoding/json"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
)

func TestTransformJSONUnmarshal(t *testing.T) {
	jsonStr := `{
		"position":[1,2,3],
		"rotation":[4,5,6,7],
		"scale":[8,9,10]
	}`

	var transform Transform
	e := json.Unmarshal([]byte(jsonStr), &transform)
	if e != nil {
		t.Fatal(e)
	}

	expectedPos := mgl.Vec3{1, 2, 3}
	if transform.position != expectedPos {
		t.Errorf("incorrect position: got %v, expected %v", transform.position, expectedPos)
	}

	expectedRot := mgl.Quat{W: 4, V: mgl.Vec3{5, 6, 7}}
	if transform.rotation != expectedRot {
		t.Errorf("incorrect rotation: got %v, expected %v", transform.rotation, expectedRot)
	}

	expectedScale := mgl.Vec3{8, 9, 10}
	if transform.scale != expectedScale {
		t.Errorf("incorrect scale: got %v, expected %v", transform.scale, expectedScale)
	}
}

func TestTransformJSONMarshal(t *testing.T) {
	expected := `{"position":[1,2,3],"rotation":[4,5,6,7],"scale":[8,9,10]}`

	transform := Transform{
		position: mgl.Vec3{1, 2, 3},
		rotation: mgl.Quat{W: 4, V: mgl.Vec3{5, 6, 7}},
		scale:    mgl.Vec3{8, 9, 10},
	}

	b, e := json.Marshal(&transform)
	if e != nil {
		t.Fatal(e)
	}

	if string(b) != expected {
		t.Errorf("JSON does not match: \ngot %v\n expected %v", string(b), expected)
	}
}

func TestTransformString(t *testing.T) {
	expected := `{"position":[1,2,3],"rotation":[4,5,6,7],"scale":[8,9,10]}`

	transform := Transform{
		position: mgl.Vec3{1, 2, 3},
		rotation: mgl.Quat{W: 4, V: mgl.Vec3{5, 6, 7}},
		scale:    mgl.Vec3{8, 9, 10},
	}

	if transform.String() != expected {
		t.Errorf("JSON does not match:\ngot %v\nexpected %v", transform.String(), expected)
	}
}

func TestTransformMatrix(t *testing.T) {
	transform := newNode()
	expected := mgl.Ident4()
	if transform.mat != expected {
		t.Errorf("incorrect matrix:\ngot %v\nexpected %v", transform.mat, expected)
	}

	transform.SetPosition(mgl.Vec3{1, 2, 3})
	expected = mgl.Translate3D(1, 2, 3)
	if transform.mat != expected {
		t.Errorf("incorrect matrix:\ngot %v\nexpected %v", transform.mat, expected)
	}
}
