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
		DiffuseTex: texture.Load("diffuse.jpg"),
		NormalTex:  texture.Load("normal.jpg"),
	}
}

func (m pbrMaterial) Apply() {
	m.Shader.Use()
	m.DiffuseTex.Bind(0)
	m.NormalTex.Bind(1)
}
