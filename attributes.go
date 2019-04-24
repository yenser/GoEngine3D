package main

import "github.com/go-gl/gl/v4.1-core/gl"

// Position Attribute
func createPositionAttribute(program uint32, name string, size, stride int32, offset int) uint32 {
	posAttrib := uint32(gl.GetAttribLocation(program, gl.Str(name)))
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, size, gl.FLOAT, false, stride*4, gl.PtrOffset(offset))
	return posAttrib
}

// Color Attribute
func createColorAttribute(program uint32, name string, size, stride int32, offset int) uint32 {
	colAttrib := uint32(gl.GetAttribLocation(program, gl.Str(name)))
	gl.EnableVertexAttribArray(colAttrib)
	gl.VertexAttribPointer(colAttrib, size, gl.FLOAT, false, stride*4, gl.PtrOffset(offset*4))
	return colAttrib
}

// Texture Attribute
func createTextureAttribute(program uint32, name string, size, stride int32, offset int) uint32 {
	texAttrib := uint32(gl.GetAttribLocation(program, gl.Str(name)))
	gl.EnableVertexAttribArray(texAttrib)
	gl.VertexAttribPointer(texAttrib, size, gl.FLOAT, false, stride*4, gl.PtrOffset(offset*4))
	return texAttrib
}
