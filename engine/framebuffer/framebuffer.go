package framebuffer

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

type FrameBuffer struct {
	handle uint32
	color  []uint32
	depth  uint32
}

func New(numCol uint32) *FrameBuffer {
	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)

	var tex = make([]uint32, numCol)
	for i := uint32(0); i < numCol; i++ {
		gl.GenTextures(1, &tex[i])
		gl.BindTexture(gl.TEXTURE_2D, tex[i])
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 1024, 768, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, gl.TEXTURE_2D, tex[i], 0)
	}

	var depth uint32
	gl.GenTextures(1, &depth)
	gl.BindTexture(gl.TEXTURE_2D, depth)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH24_STENCIL8, 1024, 768, 0, gl.DEPTH_STENCIL, gl.UNSIGNED_INT_24_8, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.TEXTURE_2D, depth, 0)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("Framebuffer not complete")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return &FrameBuffer{
		handle: fbo,
		color:  tex,
		depth:  depth,
	}
}

func Use(fbo *FrameBuffer) {
	if fbo == nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		// bufs := []uint32{gl.COLOR_ATTACHMENT0}
		// gl.DrawBuffers(1, &bufs[0])

	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, fbo.handle)
		// bufs := make([]uint32, len(fbo.color))
		// for i := 0; i < len(fbo.color); i++ {
		// 	bufs[i] = gl.COLOR_ATTACHMENT0 + uint32(i)
		// }
		// gl.DrawBuffers(int32(len(fbo.color)), &bufs[0])
	}
}

func (fbo *FrameBuffer) BindDepthTexture(target uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + target)
	gl.BindTexture(gl.TEXTURE_2D, fbo.depth)
}
func (fbo *FrameBuffer) BindColorTexture(idx uint32, target uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + target)
	gl.BindTexture(gl.TEXTURE_2D, fbo.color[idx])
}

func (fbo *FrameBuffer) Free() {
	for _, h := range fbo.color {
		gl.DeleteTextures(1, &h)
	}
	gl.DeleteTextures(1, &fbo.depth)
	gl.DeleteFramebuffers(1, &fbo.handle)
}
