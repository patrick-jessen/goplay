package texture

import (
	"image"
	"image/draw"
	_ "image/jpeg" // Support JPEG format
	_ "image/png"  // Support PNG format
	"os"

	"github.com/nfnt/resize"

	"github.com/patrick-jessen/goplay/engine/log"
	"github.com/patrick-jessen/goplay/engine/worker"

	"github.com/go-gl/gl/v3.2-core/gl"
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
		curRes:    1,
	}
)

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
	loaded  bool
	loading bool
	handle  uint32
	file    string
}

// Unload unloads the texture and its resources.
func (t *Texture) Unload() {
	gl.DeleteTextures(1, &t.handle)
	t.handle = 0
}

// Load loads the texture.
func (t *Texture) load() {
	t.loading = true
	res := Settings.curRes

	go func() {
		img := loadImage(textureDir+t.file, res)
		worker.CallSynchronized(func() {
			t.Unload()
			t.handle = newTexture(img)
			t.loading = false
			t.loaded = true
		})
	}()
}

// Bind binds the texture to the given texture location.
func (t *Texture) Bind(idx uint32) {
	if !t.loaded && !t.loading {
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
func loadImage(file string, res uint) *image.RGBA {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Panic("could not open texture file", "imgFile", imgFile, "error", err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Panic("could not decode image", "error", err)
	}

	width := img.Bounds().Dx() / int(res)
	height := img.Bounds().Dy() / int(res)
	img = resize.Resize(uint(width), uint(height), img, resize.Bicubic)

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Panic("unsupported stride", "stride", rgba.Stride)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba
}
