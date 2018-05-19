package framebuffer

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/patrick-jessen/goplay/engine/log"
)

func Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

type FrameBuffer struct {
	handle        uint32
	color         []uint32
	depth         uint32
	width, height int32
}

func New(width int, height int, numCols int, msLevel int) *FrameBuffer {
	w := int32(width)
	h := int32(height)

	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)

	var colors = make([]uint32, numCols)
	for i := uint32(0); i < uint32(numCols); i++ {
		gl.GenTextures(1, &colors[i])
		if msLevel == 0 {
			gl.BindTexture(gl.TEXTURE_2D, colors[i])
			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
			gl.BindTexture(gl.TEXTURE_2D, 0)
			gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, gl.TEXTURE_2D, colors[i], 0)
		} else {
			gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, colors[i])
			gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, int32(msLevel), gl.RGBA32F, w, h, true)
			gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, 0)
			gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+i, gl.TEXTURE_2D_MULTISAMPLE, colors[i], 0)
		}
	}

	var depth uint32
	gl.GenTextures(1, &depth)
	if msLevel == 0 {
		gl.BindTexture(gl.TEXTURE_2D, depth)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH24_STENCIL8, w, h, 0, gl.DEPTH_STENCIL, gl.UNSIGNED_INT_24_8, nil)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.TEXTURE_2D, depth, 0)
	} else {
		gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, depth)
		gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, int32(msLevel), gl.DEPTH24_STENCIL8, w, h, true)
		gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, 0)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.TEXTURE_2D_MULTISAMPLE, depth, 0)
	}

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		log.Panic("framebuffer not complete", "result", gl.CheckFramebufferStatus(gl.FRAMEBUFFER))
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return &FrameBuffer{
		handle: fbo,
		color:  colors,
		depth:  depth,
		width:  w,
		height: h,
	}
}

func (fbo *FrameBuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo.handle)
}
func (fbo *FrameBuffer) BindDepthTexture(target int) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(target))
	gl.BindTexture(gl.TEXTURE_2D, fbo.depth)
}
func (fbo *FrameBuffer) BindColorTexture(idx int, target int) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(target))
	gl.BindTexture(gl.TEXTURE_2D, uint32(fbo.color[idx]))
}
func (fbo *FrameBuffer) Blit(dst *FrameBuffer, w, h int, depth bool) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	if dst == nil {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	} else {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, dst.handle)
	}

	mask := uint32(gl.COLOR_BUFFER_BIT)
	if depth {
		mask |= gl.DEPTH_BUFFER_BIT
	}
	gl.BlitFramebuffer(0, 0, fbo.width, fbo.height, 0, 0, fbo.width, fbo.height, mask, gl.NEAREST)
}
func (fbo *FrameBuffer) Free() {
	for _, v := range fbo.color {
		gl.DeleteTextures(1, &v)
	}
	gl.DeleteTextures(1, &fbo.depth)
	gl.DeleteFramebuffers(1, &fbo.handle)
}
