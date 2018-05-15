package window

var Settings settings

type settings struct {
	vsync bool
}

func (s *settings) SetVsync(on bool) {

}
func (s *settings) Vsync() bool {
	return s.vsync
}
