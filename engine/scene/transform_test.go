package scene

import (
	"encoding/json"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl64"
)

func Test_newTransform(t *testing.T) {
	trans := newTransform()
	if trans.rotation != mgl.QuatIdent() ||
		trans.scale != (mgl.Vec3{1, 1, 1}) ||
		trans.position != (mgl.Vec3{0, 0, 0}) ||
		trans.mat != mgl.Ident4() {
		t.Errorf("wrong transformation. got %v", trans)
	}
}

func TestTransform_Position(t *testing.T) {
	expected := mgl.Vec3{1, 2, 3}
	trans := Transform{position: expected}
	if trans.Position() != expected {
		t.Errorf("wrong position. got %v, expected %v", trans.Position(), expected)
	}
}

func TestTransform_Rotation(t *testing.T) {
	expected := mgl.Quat{W: 1, V: mgl.Vec3{2, 3, 4}}
	trans := Transform{rotation: expected}
	if trans.Rotation() != expected {
		t.Errorf("wrong rotation. got %v, expected %v", trans.Rotation(), expected)
	}
}

func TestTransform_Scale(t *testing.T) {
	expected := mgl.Vec3{1, 2, 3}
	trans := Transform{scale: expected}
	if trans.Scale() != expected {
		t.Errorf("wrong scale. got %v, expected %v", trans.Scale(), expected)
	}
}

func TestTransform_update(t *testing.T) {
	trans := newTransform()

	pos := mgl.Vec3{1, 2, 3}
	rot := mgl.Quat{W: 1, V: mgl.Vec3{0, 1, 0}}
	scal := mgl.Vec3{0.1, 0.2, 0.3}

	trans.SetPosition(pos)
	trans.SetRotation(rot)
	trans.SetScale(scal)

	expected := mgl.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(rot.Mat4()).
		Mul4(mgl.Scale3D(scal.X(), scal.Y(), scal.Z()))

	if trans.mat != expected {
		t.Errorf("incorrect matrix:\ngot %v\nexpected %v", trans.mat, expected)
	}
}

func TestTransform_UnmarshalJSON(t *testing.T) {
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

func TestTransform_MarshalJSON(t *testing.T) {
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

func TestTransform_String(t *testing.T) {
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
