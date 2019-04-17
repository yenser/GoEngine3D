package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	shaders "github.com/yenser/GoEngine3D/shaders"
)

const (
	windowHeight = 600
	windowWidth  = 800
)

const SizeofFloat = 4.0

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err.Error())
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "GoEngine3D", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Important! Call gl.Init only under the presence of an active OpenGL context,
	// i.e., after MakeContextCurrent.
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(shaders.VertexShader, shaders.FragmentShader)
	if err != nil {
		panic(err)
	}

	// tell OpenGL your shader program
	gl.UseProgram(program)

	// set color variable to uniColor
	// uniColor := gl.GetUniformLocation(program, gl.Str("triangleColor\x00"))

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Declare Vertex Array Object
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Declare Vertex Buffer Object
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*SizeofFloat, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Declare Element Buffer Object
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(elements)*SizeofFloat, gl.Ptr(elements), gl.STATIC_DRAW)

	// Declare Texture
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	// (x, y, z) => (s, t, r)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &color)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 2, 2, 0, gl.RGB, gl.FLOAT, gl.Ptr(pixels))

	// Position Attribute
	posAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 2, gl.FLOAT, false, 5*SizeofFloat, gl.PtrOffset(0))

	// Color Attribute
	colAttrib := uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colAttrib)
	gl.VertexAttribPointer(colAttrib, 3, gl.FLOAT, false, 5*SizeofFloat, gl.PtrOffset(2*SizeofFloat))

	// set clear window color
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// check for errors before main window
	glErr := gl.GetError()
	if glErr != gl.NO_ERROR {
		fmt.Printf("Error: %v\n", err)
	}

	// previousTime := glfw.GetTime()

	// main run buffer
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		// time := glfw.GetTime()
		// elapsed := time - previousTime
		// previousTime = time

		// red := float32((math.Sin(time*4.0) + 1.0) / 2.0)
		// green := float32(0.0)
		// blue := float32(0.0)

		// gl.Uniform3f(uniColor, red, green, blue) // red

		// render
		gl.UseProgram(program)

		gl.BindVertexArray(vao)

		// gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

// var vertices = []float32{
// 	0.0, 0.5, 1.0, 0.0, 0.0, // Vertex 1 (X, Y)
// 	0.5, -0.5, 0.0, 1.0, 0.0, // Vertex 2 (X, Y)
// 	-0.5, -0.5, 0.0, 0.0, 1.0, // Vertex 3 (X, Y)
// }

var vertices = []float32{
	-0.5, 0.5, 1.0, 0.0, 0.0, // Top-left
	0.5, 0.5, 0.0, 1.0, 0.0, // Top-right
	0.5, -0.5, 0.0, 0.0, 1.0, // Bottom-right
	-0.5, -0.5, 1.0, 1.0, 1.0, // Bottom-left
}

var elements = []uint32{
	0, 1, 2,
	2, 3, 0,
}

var color = float32(1.0)

// var color = []float32{
// 	1.0, 0.0, 0.0, 1.0,
// }

// Black/white checkerboard
var pixels = []float32{
	0.0, 0.0, 0.0, 1.0, 1.0, 1.0,
	1.0, 1.0, 1.0, 0.0, 0.0, 0.0,
}

// var cubeVertices = []float32{
// 	//  X, Y, Z, U, V
// 	// Bottom
// 	-1.0, -1.0, -1.0, 0.0, 0.0,
// 	1.0, -1.0, -1.0, 1.0, 0.0,
// 	-1.0, -1.0, 1.0, 0.0, 1.0,
// 	1.0, -1.0, -1.0, 1.0, 0.0,
// 	1.0, -1.0, 1.0, 1.0, 1.0,
// 	-1.0, -1.0, 1.0, 0.0, 1.0,

// 	// Top
// 	-1.0, 1.0, -1.0, 0.0, 0.0,
// 	-1.0, 1.0, 1.0, 0.0, 1.0,
// 	1.0, 1.0, -1.0, 1.0, 0.0,
// 	1.0, 1.0, -1.0, 1.0, 0.0,
// 	-1.0, 1.0, 1.0, 0.0, 1.0,
// 	1.0, 1.0, 1.0, 1.0, 1.0,

// 	// Front
// 	-1.0, -1.0, 1.0, 1.0, 0.0,
// 	1.0, -1.0, 1.0, 0.0, 0.0,
// 	-1.0, 1.0, 1.0, 1.0, 1.0,
// 	1.0, -1.0, 1.0, 0.0, 0.0,
// 	1.0, 1.0, 1.0, 0.0, 1.0,
// 	-1.0, 1.0, 1.0, 1.0, 1.0,

// 	// Back
// 	-1.0, -1.0, -1.0, 0.0, 0.0,
// 	-1.0, 1.0, -1.0, 0.0, 1.0,
// 	1.0, -1.0, -1.0, 1.0, 0.0,
// 	1.0, -1.0, -1.0, 1.0, 0.0,
// 	-1.0, 1.0, -1.0, 0.0, 1.0,
// 	1.0, 1.0, -1.0, 1.0, 1.0,

// 	// Left
// 	-1.0, -1.0, 1.0, 0.0, 1.0,
// 	-1.0, 1.0, -1.0, 1.0, 0.0,
// 	-1.0, -1.0, -1.0, 0.0, 0.0,
// 	-1.0, -1.0, 1.0, 0.0, 1.0,
// 	-1.0, 1.0, 1.0, 1.0, 1.0,
// 	-1.0, 1.0, -1.0, 1.0, 0.0,

// 	// Right
// 	1.0, -1.0, 1.0, 1.0, 1.0,
// 	1.0, -1.0, -1.0, 1.0, 0.0,
// 	1.0, 1.0, -1.0, 0.0, 0.0,
// 	1.0, -1.0, 1.0, 1.0, 1.0,
// 	1.0, 1.0, -1.0, 0.0, 0.0,
// 	1.0, 1.0, 1.0, 0.0, 1.0,
// }
