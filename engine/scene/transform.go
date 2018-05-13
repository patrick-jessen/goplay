package scene

import (
	"encoding/json"

	mgl "github.com/go-gl/mathgl/mgl64"
)

// Transform represents a transformation matrix.
type Transform struct {
	position mgl.Vec3
	rotation mgl.Quat
	scale    mgl.Vec3
	mat      mgl.Mat4
}

// newTransform creates a new identity transform.
func newTransform() Transform {
	return Transform{
		rotation: mgl.QuatIdent(),
		scale:    mgl.Vec3{1, 1, 1},
		mat:      mgl.Ident4(),
	}
}

// Position returns the position.
func (t *Transform) Position() mgl.Vec3 {
	return t.position
}

// Rotation returns the rotation.
func (t *Transform) Rotation() mgl.Quat {
	return t.rotation
}

// Scale returns the scale.
func (t *Transform) Scale() mgl.Vec3 {
	return t.scale
}

// SetPosition sets the position.
func (t *Transform) SetPosition(v mgl.Vec3) {
	t.position = v
	t.update()
}

// SetRotation sets the rotation.
func (t *Transform) SetRotation(q mgl.Quat) {
	t.rotation = q
	t.update()
}

// SetScale sets the scale.
func (t *Transform) SetScale(v mgl.Vec3) {
	t.scale = v
	t.update()
}

// update updates the transformation matrix.
func (t *Transform) update() {
	t.mat = mgl.Translate3D(t.position.X(), t.position.Y(), t.position.Z())
	t.mat = t.mat.Mul4(t.rotation.Mat4())
	t.mat = t.mat.Mul4(mgl.Scale3D(t.scale.X(), t.scale.Y(), t.scale.Z()))
}

// UnmarshalJSON decodes a transform from JSON.
func (t *Transform) UnmarshalJSON(data []byte) error {
	var objMap map[string]*json.RawMessage
	e := json.Unmarshal(data, &objMap)
	if e != nil {
		return e
	}

	if pos, ok := objMap["position"]; ok {
		json.Unmarshal(*pos, &t.position)
	}
	if rot, ok := objMap["rotation"]; ok {
		var tmp [4]float64
		json.Unmarshal(*rot, &tmp)
		t.rotation = mgl.Quat{
			W: tmp[0],
			V: mgl.Vec3{
				tmp[1],
				tmp[2],
				tmp[3],
			},
		}
	}
	if scale, ok := objMap["scale"]; ok {
		json.Unmarshal(*scale, &t.scale)
	}

	t.update()
	return nil
}

// MarshalJSON encodes a transform as JSON.
func (t *Transform) MarshalJSON() ([]byte, error) {
	tmp := struct {
		P mgl.Vec3   `json:"position"`
		R [4]float64 `json:"rotation"`
		S mgl.Vec3   `json:"scale"`
	}{
		P: t.position,
		R: [4]float64{
			t.rotation.W,
			t.rotation.V[0],
			t.rotation.V[1],
			t.rotation.V[2],
		},
		S: t.scale,
	}
	return json.Marshal(&tmp)
}

// String returns a JSON encoded string.
func (t Transform) String() string {
	b, _ := json.Marshal(&t)
	return string(b)
}
