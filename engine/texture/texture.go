package texture

import (
	"image"
	"image/draw"
	_ "image/jpeg" // Support JPEG format
	_ "image/png"  // Support PNG format
	"os"
	"path/filepath"
	"strings"

	"github.com/patrick-jessen/goplay/engine/worker"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/log"
)

const (
	textureDir                 = "./assets/textures/"
	maxTextureMaxAnisotropyExt = 0x84FF
	textureMaxAnisotropyExt    = 0x84FE
)

var (
	cache    = make(map[string]*Texture)
	Settings = settings{
		curFilter: gl.LINEAR_MIPMAP_NEAREST,
		curRes:    Normal,
	}
)

type Filter int32

const (
	Bilinear  Filter = gl.LINEAR_MIPMAP_NEAREST
	Trilinear Filter = gl.LINEAR_MIPMAP_LINEAR
)

type Resolution int32

const (
	Low    Resolution = 0
	Normal Resolution = 1
	High   Resolution = 2
)

type settings struct {
	curFilter Filter
	newFilter Filter

	curAniso float32
	newAniso float32

	curRes Resolution
	newRes Resolution
}

func (s *settings) Filter() (Filter, int) {
	return s.curFilter, int(s.curAniso)
}
func (s *settings) SetFilter(f Filter, aniso int) {

	log.Info("setting filter", "f", f == Bilinear, "aniso", aniso)
	s.newFilter = f

	// Limit to max available aniso level.
	var maxAniso float32
	gl.GetFloatv(maxTextureMaxAnisotropyExt, &maxAniso)
	if int(maxAniso) < aniso {
		aniso = int(maxAniso)
	}
	s.newAniso = float32(aniso)
}

func (s *settings) Resolution() Resolution {
	return s.curRes
}
func (s *settings) SetResolution(r Resolution) {
	s.newRes = r
}

func (s *settings) Apply() {
	for _, t := range cache {

		if s.curRes != s.newRes {
			_, ok := t.fileName()
			if ok {
				// Unload texture to trigger reload of new resolution
				t.Unload()
				continue
			}
		}

		gl.BindTexture(gl.TEXTURE_2D, t.handle)
		if s.curFilter != s.newFilter {
			if s.newFilter == Bilinear {
				// Must reload to revert to bilinear
				t.Unload()
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

// Load returns a texture by either loading it or reading from cache.
func Load(name string) *Texture {
	// Read form cache
	if val, ok := cache[name]; ok {
		return val
	}
	// Load from disk
	t := Texture{file: name}
	t.load()
	cache[name] = &t
	return &t
}

// Texture represents an OpenGL texture.
type Texture struct {
	loading bool
	handle  uint32
	file    string
}

// fileName returns the actual file used (with respect to resolution).
// Second argument is true, if texture has a version for the current resolution.
func (t *Texture) fileName() (string, bool) {
	fext := filepath.Ext(t.file)
	fname := strings.TrimSuffix(t.file, fext)

	var file string
	switch Settings.curRes {
	case Low:
		file = fname + "_low" + fext
	case Normal:
		file = fname + fext
	case High:
		file = fname + "_high" + fext
	}
	if _, e := os.Stat(textureDir + file); e != nil {
		return (fname + fext), false
	}
	return file, true
}

// Unload unloads the texture and its resources.
func (t *Texture) Unload() {
	gl.DeleteTextures(1, &t.handle)
	t.handle = 0
}

// Load loads the texture.
func (t *Texture) load() {
	file, _ := t.fileName()
	t.loading = true

	go func() {
		img := loadImage(textureDir + file)
		worker.CallSynchronized(func() {
			t.handle = newTexture(img)
			t.loading = false
		})
	}()
}

// Bind binds the texture to the given texture location.
func (t *Texture) Bind(idx uint32) {
	if t.handle == 0 {
		t.load()
	}

	gl.ActiveTexture(gl.TEXTURE0 + idx)
	gl.BindTexture(gl.TEXTURE_2D, t.handle)
}

// newTexture creates and uploads the texture.
func newTexture(data *image.RGBA) uint32 {
	var handle uint32

	// Create and bind texture
	gl.GenTextures(1, &handle)
	gl.BindTexture(gl.TEXTURE_2D, handle)

	// Set parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(Settings.curFilter))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	// Anisotropic filtering
	gl.TexParameterf(gl.TEXTURE_2D, textureMaxAnisotropyExt, Settings.curAniso)

	// Upload image
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(data.Rect.Size().X),
		int32(data.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data.Pix))

	// Generate mipmaps
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return handle
}

// loadImage loads an image from file.
func loadImage(file string) *image.RGBA {
	imgFile, err := os.Open(file)
	if err != nil {
		panic("texture not found on disk")
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic("could not decode image")
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba
}
