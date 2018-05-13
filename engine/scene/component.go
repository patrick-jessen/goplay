package scene

import "reflect"

var componentMap = make(map[string]reflect.Type)

// RegisterComponent registers a component type.
func RegisterComponent(c Component) {
	typ := reflect.TypeOf(c).Elem()
	componentMap[typ.Name()] = typ
}

// Component is a behavior which can be attached to nodes.
type Component interface {
	Initialize(*Node)
	Update()
	Render()
}
