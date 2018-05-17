package material

import (
	"github.com/patrick-jessen/goplay/engine/shader"
	"github.com/patrick-jessen/goplay/engine/texture"
)

type Material struct {
	Shader   shader.Shader
	Textures []*texture.Texture
}

func New() Material {
	return Material{
		Shader: shader.Load("basic"),
		Textures: []*texture.Texture{
			texture.Load("default.png"),
		},
	}
}

func (m Material) Apply() {
	m.Shader.Use()
	for i, t := range m.Textures {
		t.Bind(uint32(i))
	}
}
