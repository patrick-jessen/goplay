package scene

import (
	"encoding/json"
	"io/ioutil"
)

const sceneDir = "./assets/scenes/"

var currentScene Scene

func Current() Scene {
	return currentScene
}

type Scene struct {
	Root *Node
}

func New() Scene {
	return Scene{
		Root: newNode(),
	}
}

func Load(name string) Scene {
	b, e := ioutil.ReadFile(sceneDir + name + ".json")
	if e != nil {
		panic("scene not found: " + e.Error())
	}

	node := newNode()
	e = json.Unmarshal(b, node)
	if e != nil {
		panic("could not unmarshal scene: " + e.Error())
	}
	node.initialize(nil, "root")
	return Scene{
		Root: node,
	}
}

func (s Scene) Update() {
	s.Root.update()
}

func (s Scene) Render() {
	s.Root.render()
}

func (s Scene) MakeCurrent() {
	currentScene = s
}
