package resource

import (
	"github.com/patrick-jessen/goplay/engine/model"
	"github.com/patrick-jessen/goplay/engine/scene"
)

func LoadScene(name string) *scene.Scene {
	s := scene.Load(name)
	for k, v := range scene.MountMap {
		model.Load(v).Mount(k)
	}
	return s
}
