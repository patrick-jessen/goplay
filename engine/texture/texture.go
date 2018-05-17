package texture

import (
	"image"
	"image/draw"
	_ "image/jpeg" // Support JPEG format
	_ "image/png"  // Support PNG format
	"os"

	"github.com/go-gl/gl/v3.2-core/gl"
)

const textureDir = "./assets/textures/"
const (
	maxTextureMaxAnisotropyExt = 0x84FF
	textureMaxAnisotropyExt    = 0x84FE
)

var cache = make(map[string]Texture)

// Load returns a texture by either loading it or reading from cache.
func Load(name string) Texture {
	// Read form cache
	if val, ok := cache[name]; ok {
		return val
	}
	// Load from disk
	cache[name] = Texture{handle: loadTexture(name)}
	return cache[name]
}

// Texture represents an OpenGL texture.
type Texture struct {
	handle uint32
}

// Bind binds the texture to the given texture location.
func (t Texture) Bind(idx uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + idx)
	gl.BindTexture(gl.TEXTURE_2D, t.handle)
}

// loadTexture loads a texture from a file.
func loadTexture(name string) uint32 {
	file := textureDir + name
	if _, err := os.Stat(file); err == nil {

	} else if _, err := os.Stat(file + ".png"); err == nil {
		file += ".png"
	} else if _, err := os.Stat(file + ".jpg"); err == nil {
		file += ".jpg"
	} else {
		panic("Texture not found: " + name)
	}

	img := loadImage(file)
	return newTexture(img)
}

// newTexture creates and uploads the texture.
func newTexture(data *image.RGBA) uint32 {
	var handle uint32

	// Create and bind texture
	gl.GenTextures(1, &handle)
	gl.BindTexture(gl.TEXTURE_2D, handle)

	// Set parameters
	var minFilter int32 = gl.NEAREST
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, minFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	// Anisotropic filtering
	var aniso float32
	gl.GetFloatv(maxTextureMaxAnisotropyExt, &aniso)
	gl.TexParameterf(gl.TEXTURE_2D, textureMaxAnisotropyExt, aniso)

	// Upload image
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(data.Rect.Size().X),
		int32(data.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data.Pix))

	// Generate mipmap if required
	if minFilter == gl.NEAREST_MIPMAP_NEAREST ||
		minFilter == gl.LINEAR_MIPMAP_NEAREST ||
		minFilter == gl.NEAREST_MIPMAP_LINEAR ||
		minFilter == gl.LINEAR_MIPMAP_LINEAR {
		gl.GenerateMipmap(gl.TEXTURE_2D)
	}

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
