package texture

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/log"
)

var Settings = settings{
	curFilter: Trilinear,
	newFilter: Trilinear,
	curRes:    1,
	newRes:    1,
	curAniso:  16,
	newAniso:  16,
}

type Filter int32

const (
	Bilinear  Filter = gl.LINEAR_MIPMAP_NEAREST
	Trilinear Filter = gl.LINEAR_MIPMAP_LINEAR
)

type settings struct {
	curFilter Filter
	newFilter Filter

	curAniso float32
	newAniso float32

	curRes uint
	newRes uint
}

func (s *settings) Filter() (Filter, int) {
	return s.curFilter, int(s.curAniso)
}
func (s *settings) SetFilter(f Filter, aniso int) {
	s.newFilter = f
	s.newAniso = float32(aniso)

	// Limit to max available aniso level.
	var maxAniso float32
	gl.GetFloatv(maxTextureMaxAnisotropyExt, &maxAniso)
	if maxAniso < s.newAniso {
		s.newAniso = maxAniso
	}
}

func (s *settings) Resolution() uint {
	return s.curRes
}
func (s *settings) SetResolution(r uint) {
	if r == 0 {
		log.Warn("resolution must be at least 1", "r", r)
		r = 1
	}
	s.newRes = r
}

func (s *settings) Apply() {
	for _, t := range cache {

		if s.curRes != s.newRes {
			// Unload texture to trigger reload of new resolution
			t.loaded = false
			continue
		}

		gl.BindTexture(gl.TEXTURE_2D, t.handle)
		if s.curFilter != s.newFilter {
			if s.newFilter == Bilinear {
				// Must reload to revert to bilinear
				t.loaded = false
				continue
			}
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(s.newFilter))
		}
		if s.curAniso != s.newAniso {
			gl.TexParameterf(gl.TEXTURE_2D, textureMaxAnisotropyExt, s.newAniso)
		}
	}

	s.curFilter = s.newFilter
	s.curAniso = s.newAniso
	s.curRes = s.newRes
}
