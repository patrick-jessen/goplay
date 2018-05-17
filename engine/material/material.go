package material

import (
	"github.com/patrick-jessen/goplay/engine/shader"
	"github.com/patrick-jessen/goplay/engine/texture"
)

type Material interface {
	Apply()
}

type pbrMaterial struct {
	Shader     shader.Shader
	DiffuseTex *texture.Texture
	NormalTex  *texture.Texture
}

func NewPBRMaterial() pbrMaterial {
	return pbrMaterial{
		Shader:     shader.Load("pbr"),
		DiffuseTex: texture.Load("default_diff.jpg"),
		NormalTex:  texture.Load("default_norm.jpg"),
	}
}

func (m pbrMaterial) Apply() {
	m.Shader.Use()
	m.DiffuseTex.Bind(0)
	m.NormalTex.Bind(1)
}

type defaultMaterial struct{}

func NewDefaultMaterial() Material {
	return defaultMaterial{}
}
func (defaultMaterial) Apply() {}
