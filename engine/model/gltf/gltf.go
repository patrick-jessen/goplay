// Package gltf implements a glTF 2.0 loader.
// The package does NOT peform any validation.
package gltf

import "encoding/json"

////////////////////////////////////////////////////////////////////////////////
// Accessor
////////////////////////////////////////////////////////////////////////////////

// Accessor is a typed view into a bufferView.
type Accessor struct {
	BufferView    uint           `json:"bufferView"`    // The index of the bufferView.
	ByteOffset    uint           `json:"byteOffset"`    // The offset relative to the start of the bufferView in bytes.
	ComponentType uint           `json:"componentType"` // The datatype of components in the attribute.
	Normalized    bool           `json:"normalized"`    // Specifies whether integer data values should be normalized.
	Count         uint           `json:"count"`         // The number of attributes referenced by this accessor.
	Type          string         `json:"type"`          // Specifies if the attribute is a scalar, vector, or matrix.
	Max           []uint         `json:"max"`           // Maximum value of each component in this attribute.
	Min           []uint         `json:"min"`           // Minimum value of each component in this attribute.
	Name          string         `json:"name"`          // The name of the accessor
	Sparse        AccessorSparse `json:"sparse"`        // Sparse storage of attributes that deviate from their initialization value.
}

// AccessorSparse is the sparse storage of attributes that deviate from their initialization value.
type AccessorSparse struct {
	Count   uint                  `json:"count"`   // Number of entries stored in the sparse array.
	Indices AccessorSparseIndices `json:"indices"` // Array of indices that points to the respective accessor attributes.
	Values  AccessorSparseValues  `json:"values"`  // Stores the displaced accessor attributes.
}

// AccessorSparseIndices are indices of those attributes that deviate from their initialization value.
type AccessorSparseIndices struct {
	BufferView    uint `json:"bufferView"`    // The index of the bufferView with sparse indices.
	ByteOffset    uint `json:"byteOffset"`    // The offset relative to the start of the bufferView in bytes.
	ComponentType uint `json:"componentType"` // The indices data type.
}

// AccessorSparseValues stores the displaced accessor attributes.
type AccessorSparseValues struct {
	BufferView uint `json:"bufferView"` // The index of the bufferView with sparse values.
	ByteOffset uint `json:"byteOffset"` // The offset relative to the start of the bufferView in bytes.
}

////////////////////////////////////////////////////////////////////////////////
// Animation
////////////////////////////////////////////////////////////////////////////////

// Animation is a keyframe animation.
type Animation struct {
	Channels []AnimationChannel `json:"channels"` // An array of channels, each of which targets an animation's sampler at a node's property.
	Samplers []AnimationSampler `json:"samplers"` // An array of samplers that combines input and output accessors with an interpolation algorithm.
	Name     string             `json:"name"`     // The name of the animation.
}

// AnimationChannel targets an animation's sampler at a node's property.
type AnimationChannel struct {
	Sampler uint                   `json:"sampler"` // The index of a sampler in this animation used to compute the value for the target.
	Target  AnimationChannelTarget `json:"target"`  // The index of the node and TRS property to target.
}

// AnimationChannelTarget is the index of the node and TRS property that an animation channel targets.
type AnimationChannelTarget struct {
	Node uint   `json:"node"` // The index of the node to target.
	Path string `json:"path"` // The name of the node's TRS property to modify.
}

// AnimationSampler combines input and output accessors with an interpolation algorithm to define a keyframe graph.
type AnimationSampler struct {
	Input         uint   `json:"input"`         // The index of an accessor containing keyframe input values.
	Interpolation string `json:"interpolation"` // Interpolation algorithm.
	Output        uint   `json:"output"`        // The index of an accessor, containing keyframe output values.
}

// UnmarshalJSON sets default values for AnimationSampler.
func (as *AnimationSampler) UnmarshalJSON(d []byte) error {
	type alias AnimationSampler
	out := &alias{
		Interpolation: "LINEAR",
	}
	e := json.Unmarshal(d, out)
	*as = AnimationSampler(*out)
	return e
}

////////////////////////////////////////////////////////////////////////////////
// Asset
////////////////////////////////////////////////////////////////////////////////

// Asset holds metadata about the glTF asset.
type Asset struct {
	Copyright  string `json:"copyright"`  // A copyright message suitable for display to credit the content creator.
	Generator  string `json:"generator"`  // Tool that generated this glTF model.
	Version    string `json:"version"`    // The glTF version that this asset targets.
	MinVersion string `json:"minVersion"` // The minimum glTF version that this asset targets.
}

////////////////////////////////////////////////////////////////////////////////
// Buffer
////////////////////////////////////////////////////////////////////////////////

// Buffer points to binary geometry, animation, or skins.
type Buffer struct {
	URI        string `json:"uri"`        // The uri of the buffer.
	ByteLength uint   `json:"byteLength"` // The length of the buffer in bytes.
	Name       string `json:"name"`       // The name of the buffer
}

////////////////////////////////////////////////////////////////////////////////
// BufferView
////////////////////////////////////////////////////////////////////////////////

// BufferView is a view into a buffer generally representing a subset of the buffer
type BufferView struct {
	Buffer     uint   `json:"buffer"`     // The index of the buffer.
	ByteOffset uint   `json:"byteOffset"` // The offset into the buffer in bytes.
	ByteLength uint   `json:"byteLength"` // The length of the bufferView in bytes.
	ByteStride uint   `json:"byteStride"` // The stride, in bytes.
	Target     uint   `json:"target"`     // The target that the GPU buffer should be bound to.
	Name       string `json:"name"`       // The name of the buffer view
}

////////////////////////////////////////////////////////////////////////////////
// Camera
////////////////////////////////////////////////////////////////////////////////

// Camera holds a camera's projection.
type Camera struct {
	Orthographic CameraOrthographic `json:"orthographic"` // Properties for creating an orthographic projection matrix.
	Perspective  CameraPerspective  `json:"perspective"`  // Properties for creating a perspective projection matrix.
	Type         string             `json:"type"`         // Specifies if the camera uses a perspective or orthographic projection.
	Name         string             `json:"name"`         // The name of the camera.
}

// CameraOrthographic is an orthographic camera containing properties to create an orthographic projection matrix.
type CameraOrthographic struct {
	Xmag  float32 `json:"xmag"`  // The floating-point horizontal magnification of the view.
	Ymag  float32 `json:"ymag"`  // The floating-point vertical magnification of the view.
	Zfar  float32 `json:"zfar"`  // The floating-point distance to the far clipping plane.
	Znear float32 `json:"znear"` // The floating-point distance to the near clipping plane.
}

// CameraPerspective is a perspective camera containing properties to create a perspective projection matrix.
type CameraPerspective struct {
	AspectRatio float32 `json:"aspectRatio"` // The floating-point aspect ratio of the field of view.
	Yfov        float32 `json"yfov"`         // The floating-point vertical field of view in radians.
	Zfar        float32 `json:"zfar"`        // The floating-point distance to the far clipping plane.
	Znear       float32 `json:"znear"`       // The floating-point distance to the near clipping plane.
}

////////////////////////////////////////////////////////////////////////////////
// glTF
////////////////////////////////////////////////////////////////////////////////

// GlTF is the root object for a glTF asset
type GlTF struct {
	ExtensionsUsed     []string     `json:"extensionsUsed"`     // Names of glTF extensions used somewhere in this asset.
	ExtensionsRequired []string     `json:"extensionsRequired"` // Names of glTF extensions required to properly load this asset.
	Accessors          []Accessor   `json:"accessors"`          // An array of accessors.
	Animations         []Animation  `json"animations"`          // An array of keyframe animations.
	Asset              Asset        `json:"asset"`              // Metadata about the glTF asset.
	Buffers            []Buffer     `json:"buffers"`            // An array of buffers.
	BufferViews        []BufferView `json:"bufferViews"`        // An array of bufferViews.
	Cameras            []Camera     `json:"cameras"`            // An array of cameras.
	Images             []Image      `json:"images"`             // An array of images.
	Materials          []Material   `json:"materials"`          // An array of materials.
	Meshes             []Mesh       `json:"meshes"`             // An array of meshes.
	Nodes              []Node       `json:"nodes"`              // An array of nodes.
	Samplers           []Sampler    `json:"samplers"`           // An array of samplers.
	Scene              uint         `json:"scene"`              // The index of the default scene.
	Scenes             []Scene      `json:"scenes"`             // An array of scenes.
	Skins              []Skin       `json:"skins"`              // An array of skins.
	Textures           []Texture    `json:"textures"`           // An array of textures.
}

////////////////////////////////////////////////////////////////////////////////
// Image
////////////////////////////////////////////////////////////////////////////////

// Image holds data used to create a texture.
type Image struct {
	URI        string `json:"uri"`        // The uri of the image.
	MimeType   string `json:"mimeType"`   // The image's MIME type.
	BufferView uint   `json:"bufferView"` // The index of the bufferView that contains the image.
	Name       string `json:"name"`       // The name of the image.
}

////////////////////////////////////////////////////////////////////////////////
// Material
////////////////////////////////////////////////////////////////////////////////

// Material describes the material appearance of a primitive.
type Material struct {
	Name                 string                       `json:"name"`                 // The name of the material.
	PbrMetallicRoughness MaterialPbrMetallicRoughness `json:"pbrMetallicRoughness"` // A set of parameters used to define the metallic-roughness material model.
	NormalTexture        MaterialNormalTextureInfo    `json:"normalTexture"`        // A tangent space normal map.
	OcclusionTexture     MaterialOcclusionTextureInfo `json:"occlusionTexture"`     // The occlusion map texture.
	EmissiveTexture      TextureInfo                  `json:"emissiveTexture"`      // The emissive map texture.
	EmissiveFactor       []float32                    `json:"emissiveFactor"`       // The emissive color of the material.
	AlphaMode            string                       `json:"alphaMode"`            // The alpha rendering mode of the material.
	AlphaCutoff          float32                      `json:"alphaCutoff"`          // The alpha cutoff value of the material.
	DoubleSided          bool                         `json:"doubleSided"`          // Specifies whether the material is double sided.
}

// UnmarshalJSON sets default values for Material.
func (m *Material) UnmarshalJSON(d []byte) error {
	type alias Material
	out := &alias{
		EmissiveFactor: []float32{0, 0, 0},
		AlphaMode:      "OPAQUE",
		AlphaCutoff:    0.5,
		NormalTexture: MaterialNormalTextureInfo{
			Scale:    1,
			Index:    -1,
			TexCoord: -1,
		},
	}
	e := json.Unmarshal(d, out)
	*m = Material(*out)
	return e
}

// MaterialPbrMetallicRoughness describes a material PBR Metallic Roughness.
type MaterialPbrMetallicRoughness struct {
	BaseColorFactor          []float32   `json:"baseColorFactor"`
	BaseColorTexture         TextureInfo `json:"baseColorTexture"`
	MetallicFactor           float32     `json:"metallicFactor"`
	RoughnessFactor          float32     `json:"roughnessFactor"`
	MetallicRoughnessTexture TextureInfo `json:"metallicRoughnessTexture"`
}

// UnmarshalJSON sets default values for MaterialPbrMetallicRoughness.
func (m *MaterialPbrMetallicRoughness) UnmarshalJSON(d []byte) error {
	type alias MaterialPbrMetallicRoughness
	out := &alias{
		BaseColorFactor: []float32{1, 1, 1, 1},
		BaseColorTexture: TextureInfo{
			Index:    -1,
			TexCoord: -1,
		},
		MetallicFactor:  1,
		RoughnessFactor: 1,
	}
	e := json.Unmarshal(d, out)
	*m = MaterialPbrMetallicRoughness(*out)
	return e
}

// MaterialNormalTextureInfo holds material normal texture info.
type MaterialNormalTextureInfo struct {
	Scale    float32 `json:"scale"` // The scalar multiplier applied to each normal vector of the normal texture.
	Index    int     `json:"index"`
	TexCoord int     `json:"texCoord"`
}

// UnmarshalJSON sets default values for MaterialNormalTextureInfo.
// func (m *MaterialNormalTextureInfo) UnmarshalJSON(d []byte) error {
// 	type alias MaterialNormalTextureInfo
// 	out := &alias{
// 		Scale:    1,
// 		Index:    -1,
// 		TexCoord: -1,
// 	}
// 	e := json.Unmarshal(d, out)
// 	*m = MaterialNormalTextureInfo(*out)
// 	return e
// }

// MaterialOcclusionTextureInfo holds material occlusion texture info
type MaterialOcclusionTextureInfo struct {
	Strength float32 `json:"strength"` // A scalar multiplier controlling the amount of occlusion applied.
	// Ignored properties:
	// * index
	// * texCoord
}

// UnmarshalJSON sets default values for MaterialOcclusionTextureInfo.
func (m *MaterialOcclusionTextureInfo) UnmarshalJSON(d []byte) error {
	type alias MaterialOcclusionTextureInfo
	out := &alias{
		Strength: 1,
	}
	e := json.Unmarshal(d, out)
	*m = MaterialOcclusionTextureInfo(*out)
	return e
}

////////////////////////////////////////////////////////////////////////////////
// Mesh
////////////////////////////////////////////////////////////////////////////////

// Mesh is a set of primitives to be rendered.
type Mesh struct {
	Primitives []MeshPrimitive `json:"primitives"` // An array of primitives, each defining geometry to be rendered with a material.
	Weights    []float32       `json:"weights"`    // Array of weights to be applied to the Morph Targets.
	Name       string          `json:"name"`       // The name of the mesh
}

// MeshPrimitive is the geometry to be rendered with the given material.
type MeshPrimitive struct {
	Attributes map[string]uint   `json:"attributes"` // A dictionary object, where each key corresponds to mesh attribute semantic.
	Indices    int               `json:"indices"`    // The index of the accessor that contains the indices.
	Material   int               `json:"material"`   // The index of the material to apply to this primitive when rendering.
	Mode       uint              `json:"mode"`       // The type of primitives to render.
	Targets    []map[string]uint `json:"targets"`    // An array of Morph Targets
}

// UnmarshalJSON sets default values for MeshPrimitive.
func (m *MeshPrimitive) UnmarshalJSON(d []byte) error {
	type alias MeshPrimitive
	out := &alias{
		Indices:  -1,
		Material: -1,
		Mode:     4,
	}
	e := json.Unmarshal(d, out)
	*m = MeshPrimitive(*out)
	return e
}

////////////////////////////////////////////////////////////////////////////////
// Node
////////////////////////////////////////////////////////////////////////////////

// Node is a node in the node hierarchy.
type Node struct {
	Camera      uint      `json:"camera"`      // The index of the camera referenced by this node.
	Children    []uint    `json:"children"`    // The indices of this node's children.
	Skin        uint      `json:"skin"`        // The index of the skin referenced by this node.
	Matrix      []float32 `json:"matrix"`      // A floating-point 4x4 transformation matrix stored in column-major order.
	Mesh        int       `json:"mesh"`        // The index of the mesh in this node.
	Rotation    []float32 `json:"rotation"`    // The node's unit quaternion rotation.
	Scale       []float32 `json:"scale"`       // The node's non-uniform scale
	Translation []float32 `json:"translation"` // The node's translation along the x, y, and z axes.
	Weights     []float32 `json:"weights"`     // The weights of the instantiated Morph Target.
	Name        string    `json:"name"`        // The name of the node
}

// UnmarshalJSON sets default values for Node.
func (n *Node) UnmarshalJSON(d []byte) error {
	type alias Node
	out := &alias{
		Mesh:        -1,
		Matrix:      []float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
		Rotation:    []float32{0, 0, 0, 1},
		Scale:       []float32{1, 1, 1},
		Translation: []float32{0, 0, 0},
	}
	e := json.Unmarshal(d, out)
	*n = Node(*out)
	return e
}

////////////////////////////////////////////////////////////////////////////////
// Sampler
////////////////////////////////////////////////////////////////////////////////

// Sampler holds texture sampler properties for filtering and wrapping modes.
type Sampler struct {
	MagFilter uint   `json:"magFilter"` // Magnification filter.
	MinFilter uint   `json:"minFilter"` // Minification filter.
	WrapS     uint   `json:"wrapS"`     // s wrapping mode.
	WrapT     uint   `json:"wrapT"`     // t wrapping mode.
	Name      string `json:"name"`      // The name of the sampler.
}

// UnmarshalJSON sets default values for Sampler.
func (s *Sampler) UnmarshalJSON(d []byte) error {
	type alias Sampler
	out := &alias{
		WrapS: 10497,
		WrapT: 10497,
	}
	e := json.Unmarshal(d, out)
	*s = Sampler(*out)
	return e
}

////////////////////////////////////////////////////////////////////////////////
// Scene
////////////////////////////////////////////////////////////////////////////////

// Scene contains the root nodes of a scene.
type Scene struct {
	Nodes []uint `json:"nodes"` // The indices of each root node.
	Name  string `json:"name"`  // The name of the scene.
}

////////////////////////////////////////////////////////////////////////////////
// Skin
////////////////////////////////////////////////////////////////////////////////

// Skin holds joints and matrices defining a skin.
type Skin struct {
	InverseBindMatrices uint   `json:"inverseBindMatrices"` // The index of the accessor containing the floating-point 4x4 inverse-bind matrices.
	Skeleton            uint   `json:"skeleton"`            // The index of the node used as a skeleton root.
	Joints              []uint `json:"joints"`              // Indices of skeleton nodes, used as joints in this skin.
	Name                string `json:"name"`                // The name of the skin.
}

////////////////////////////////////////////////////////////////////////////////
// Texture
////////////////////////////////////////////////////////////////////////////////

// Texture holds a texture and its sampler.
type Texture struct {
	Sampler uint   `json:"sampler"` // The index of the sampler used by this texture.
	Source  uint   `json:"source"`  // The index of the image used by this texture.
	Name    string `json:"name"`    // The name of the texture.
}

////////////////////////////////////////////////////////////////////////////////
// TextureInfo
////////////////////////////////////////////////////////////////////////////////

// TextureInfo holds a reference to a texture.
type TextureInfo struct {
	Index    int `json:"index"`    // The index of the texture.
	TexCoord int `json:"texCoord"` // The set index of texture's TEXCOORD attribute used for texture coordinate mapping.
}
