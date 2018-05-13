package scene

type Scene struct {
	Root *Node
}

func New() Scene {
	return Scene{
		Root: newNode(),
	}
}

func (s Scene) Update() {
	s.Root.update()
}

func (s Scene) Render() {
	s.Root.render()
}
