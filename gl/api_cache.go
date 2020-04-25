package gl

import (
	"unsafe"
)

type apiCache struct {
	api                API
	framebuffer        *uint32
	viewport           *viewport
	scissorTestEnabled *bool
	blendEnabled       *bool
}

type viewport struct {
	x, y, w, h int32
}

func (a *apiCache) GenBuffers(n int32, buffers *uint32) {
	a.api.GenBuffers(n, buffers)
}

func (a *apiCache) BindBuffer(target uint32, buffer uint32) {
	a.api.BindBuffer(target, buffer)
}

func (a *apiCache) BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
	a.api.BufferData(target, size, data, usage)
}

func (a *apiCache) BufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {
	a.api.BufferSubData(target, offset, size, data)
}

func (a *apiCache) GetBufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {
	a.api.GetBufferSubData(target, offset, size, data)
}

func (a *apiCache) DeleteBuffers(n int32, buffers *uint32) {
	a.api.DeleteBuffers(n, buffers)
}

func (a *apiCache) GenVertexArrays(n int32, arrays *uint32) {
	a.api.GenVertexArrays(n, arrays)
}

func (a *apiCache) DeleteVertexArrays(n int32, arrays *uint32) {
	a.api.DeleteVertexArrays(n, arrays)
}

func (a *apiCache) BindVertexArray(array uint32) {
	a.api.BindVertexArray(array)
}

func (a *apiCache) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	a.api.VertexAttribPointer(index, size, xtype, normalized, stride, pointer)
}

func (a *apiCache) EnableVertexAttribArray(index uint32) {
	a.api.EnableVertexAttribArray(index)
}

func (a *apiCache) CreateShader(xtype uint32) uint32 {
	return a.api.CreateShader(xtype)
}

func (a *apiCache) ShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {
	a.api.ShaderSource(shader, count, xstring, length)

}

func (a *apiCache) CompileShader(shader uint32) {
	a.api.CompileShader(shader)
}

func (a *apiCache) GetShaderiv(shader uint32, pname uint32, params *int32) {
	a.api.GetShaderiv(shader, pname, params)
}

func (a *apiCache) GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
	a.api.GetShaderInfoLog(shader, bufSize, length, infoLog)
}

func (a *apiCache) DeleteShader(shader uint32) {
	a.api.DeleteShader(shader)
}

func (a *apiCache) AttachShader(program uint32, shader uint32) {
	a.api.AttachShader(program, shader)
}

func (a *apiCache) LinkProgram(program uint32) {
	a.api.LinkProgram(program)
}

func (a *apiCache) GetProgramiv(program uint32, pname uint32, params *int32) {
	a.api.GetProgramiv(program, pname, params)
}

func (a *apiCache) GetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *uint8) {
	panic("implement me")
}

func (a *apiCache) UseProgram(program uint32) {
	a.api.UseProgram(program)
}

func (a *apiCache) CreateProgram() uint32 {
	return a.api.CreateProgram()
}

func (a *apiCache) GetActiveUniform(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
	a.api.GetActiveUniform(program, index, bufSize, length, size, xtype, name)
}

func (a *apiCache) GetActiveAttrib(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
	a.api.GetActiveAttrib(program, index, bufSize, length, size, xtype, name)
}

func (a *apiCache) GetAttribLocation(program uint32, name *uint8) int32 {
	return a.api.GetAttribLocation(program, name)
}

var trueVal = true

func (a *apiCache) Enable(cap uint32) {
	if cap == scissorTest {
		if a.scissorTestEnabled != nil && *a.scissorTestEnabled {
			return
		}
		a.scissorTestEnabled = &trueVal
	}
	if cap == blend {
		if a.blendEnabled != nil && *a.blendEnabled {
			return
		}
		a.blendEnabled = &trueVal
	}
	a.api.Enable(cap)
}

func (a *apiCache) BindFramebuffer(target uint32, framebufferName uint32) {
	if target == framebuffer {
		if a.framebuffer != nil && *a.framebuffer == framebufferName {
			return
		}
		a.framebuffer = &framebufferName
	}
	a.api.BindFramebuffer(target, framebufferName)
}

func (a *apiCache) Scissor(x int32, y int32, width int32, height int32) {
	a.api.Scissor(x, y, width, height)
}

func (a *apiCache) Viewport(x int32, y int32, width int32, height int32) {
	if a.viewport != nil {
		cachedViewport := *a.viewport
		if cachedViewport.x == x && cachedViewport.y == y && cachedViewport.w == width && cachedViewport.h == height {
			return
		}
		a.viewport.x = x
		a.viewport.y = y
		a.viewport.w = width
		a.viewport.h = height
	} else {
		a.viewport = &viewport{
			x: x,
			y: y,
			w: width,
			h: height,
		}
	}
	a.api.Viewport(x, y, width, height)
}

func (a *apiCache) ClearColor(red float32, green float32, blue float32, alpha float32) {
	a.api.ClearColor(red, green, blue, alpha)
}

func (a *apiCache) Clear(mask uint32) {
	a.api.Clear(mask)
}

func (a *apiCache) DrawArrays(mode uint32, first int32, count int32) {
	a.api.DrawArrays(mode, first, count)
}

func (a *apiCache) Uniform1f(location int32, v0 float32) {
	panic("implement me")
}

func (a *apiCache) Uniform2f(location int32, v0 float32, v1 float32) {
	panic("implement me")
}

func (a *apiCache) Uniform3f(location int32, v0 float32, v1 float32, v2 float32) {
	panic("implement me")
}

func (a *apiCache) Uniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32) {
	panic("implement me")
}

func (a *apiCache) Uniform1i(location int32, v0 int32) {
	a.api.Uniform1i(location, v0)
}

func (a *apiCache) Uniform2i(location int32, v0 int32, v1 int32) {
	panic("implement me")
}

func (a *apiCache) Uniform3i(location int32, v0 int32, v1 int32, v2 int32) {
	panic("implement me")
}

func (a *apiCache) Uniform4i(location int32, v0 int32, v1 int32, v2 int32, v3 int32) {
	panic("implement me")
}

func (a *apiCache) UniformMatrix3fv(location int32, count int32, transpose bool, value *float32) {
	panic("implement me")
}

func (a *apiCache) UniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {
	panic("implement me")
}

func (a *apiCache) ActiveTexture(texture uint32) {
	a.api.ActiveTexture(texture)
}

func (a *apiCache) BindTexture(target uint32, texture uint32) {
	a.api.BindTexture(target, texture)
}

func (a *apiCache) GetIntegerv(pname uint32, data *int32) {
	panic("implement me")
}

func (a *apiCache) GenTextures(n int32, textures *uint32) {
	a.api.GenTextures(n, textures)
}

func (a *apiCache) TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	a.api.TexImage2D(target, level, internalformat, width, height, border, format, xtype, pixels)
}

func (a *apiCache) TexParameteri(target uint32, pname uint32, param int32) {
	a.api.TexParameteri(target, pname, param)
}

func (a *apiCache) GenFramebuffers(n int32, framebuffers *uint32) {
	a.api.GenFramebuffers(n, framebuffers)
}

func (a *apiCache) FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	a.api.FramebufferTexture2D(target, attachment, textarget, texture, level)
}

func (a *apiCache) TexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	a.api.TexSubImage2D(target, level, xoffset, yoffset, width, height, format, xtype, pixels)
}

func (a *apiCache) GetTexImage(target uint32, level int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	panic("implement me")
}

func (a *apiCache) GetError() uint32 {
	panic("implement me")
}

func (a *apiCache) ReadPixels(x int32, y int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	panic("implement me")
}

func (a *apiCache) BlendFunc(sfactor uint32, dfactor uint32) {
	a.api.BlendFunc(sfactor, dfactor)
}

func (a *apiCache) Finish() {
	a.api.Finish()
}

func (a *apiCache) Ptr(data interface{}) unsafe.Pointer {
	return a.api.Ptr(data)
}

func (a *apiCache) PtrOffset(offset int) unsafe.Pointer {
	return a.api.PtrOffset(offset)
}

func (a *apiCache) GoStr(cstr *uint8) string {
	return a.api.GoStr(cstr)
}

func (a *apiCache) Strs(strs ...string) (cstrs **uint8, free func()) {
	return a.api.Strs(strs...)
}
