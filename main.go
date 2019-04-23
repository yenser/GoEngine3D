package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	glm "github.com/go-gl/mathgl/mgl32"
	shaders "github.com/yenser/GoEngine3D/shaders"
)

const (
	windowHeight = 600
	windowWidth  = 800
	windowTitle  = "GoEngine3D"
	enableMSAA   = true
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
	if enableMSAA {
		glfw.WindowHint(glfw.Samples, 4)
	}

	window, err := glfw.CreateWindow(windowWidth, windowHeight, windowTitle, nil, nil)
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
	uniColor := gl.GetUniformLocation(program, gl.Str("overrideColor\x00"))

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 1.0, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 1})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// bind Frag Data Location
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
	// var ebo uint32
	// gl.GenBuffers(1, &ebo)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(elements)*SizeofFloat, gl.Ptr(elements), gl.STATIC_DRAW)

	// Declare Texture
	texture, err := newTexture("./textures/diamonds.png")
	// texture, err := newTexture("./textures/texture.jpg")
	if err != nil {
		log.Fatalln(err)
	}

	// Position Attribute
	posAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 3, gl.FLOAT, false, 8*SizeofFloat, gl.PtrOffset(0))

	// Color Attribute
	colAttrib := uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colAttrib)
	gl.VertexAttribPointer(colAttrib, 3, gl.FLOAT, false, 8*SizeofFloat, gl.PtrOffset(3*SizeofFloat))

	// Texture Attribute
	texAttrib := uint32(gl.GetAttribLocation(program, gl.Str("texcoord\x00")))
	gl.EnableVertexAttribArray(texAttrib)
	gl.VertexAttribPointer(texAttrib, 2, gl.FLOAT, false, 8*SizeofFloat, gl.PtrOffset(6*SizeofFloat))

	// set clear window color
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	// set multisampling
	if enableMSAA {
		gl.Enable(gl.MULTISAMPLE)
	}

	// check for errors before main window
	glErr := gl.GetError()
	fmt.Printf("Error Code: %v\n", glErr)

	angle := 0.0
	previousTime := glfw.GetTime()
	lastTime := previousTime
	nbFrames := 0

	gl.Enable(gl.DEPTH_TEST) // enable Z-buffer

	// main run buffer
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed

		nbFrames++
		if (time - lastTime) >= 1.0 {
			title := fmt.Sprintf("%v %v ms/frame", windowTitle, 1000.0/nbFrames)
			window.SetTitle(title)
			nbFrames = 0
			lastTime += 1.0
		}

		// render
		gl.UseProgram(program)

		// Draw cube and floor
		model = glm.HomogRotate3D(float32(angle), glm.Vec3{0, 0, 1})
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		gl.Enable(gl.STENCIL_TEST)

		// Draw Floor
		gl.StencilFunc(gl.ALWAYS, 1, 0xFF)
		gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE)
		gl.StencilMask(0xFF)
		gl.DepthMask(false)
		gl.Clear(gl.STENCIL_BUFFER_BIT)

		gl.DrawArrays(gl.TRIANGLES, 36, 6)

		// Draw reflection
		gl.StencilFunc(gl.EQUAL, 1, 0xFF)
		gl.StencilMask(0x00)
		gl.DepthMask(true)

		model = model.Mul4(glm.Translate3D(0, 0, -1))
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.Uniform3f(uniColor, 0.3, 0.3, 0.3)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		gl.Uniform3f(uniColor, 1, 1, 1)

		gl.Disable(gl.STENCIL_TEST)

		// bind vertex array
		gl.BindVertexArray(vao)

		// activate and bind texture
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}

}

// var vertices = []float32{
// 	// X    Y    Z    R    G    B    U    V
// 	-0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0,
// 	0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0,
// 	0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 1.0, 1.0,
// 	-0.5, -0.5, 0.0, 1.0, 1.0, 1.0, 0.0, 1.0,
// }

var vertices = []float32{
	//X     Y     Z    R    G    B    U    V
	-0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0, 1.0, 1.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0, 1.0, 1.0, 1.0,
	-0.5, 0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0,
	-0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 0.0,

	-0.5, -0.5, 0.5, 1.0, 1.0, 1.0, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 1.0,
	-0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 0.0, 1.0,
	-0.5, -0.5, 0.5, 1.0, 1.0, 1.0, 0.0, 0.0,

	-0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0,
	-0.5, 0.5, -0.5, 1.0, 1.0, 1.0, 1.0, 1.0,
	-0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0,
	-0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0,
	-0.5, -0.5, 0.5, 1.0, 1.0, 1.0, 0.0, 0.0,
	-0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0, 1.0, 1.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 1.0, 1.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0,

	-0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 1.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0,
	-0.5, -0.5, 0.5, 1.0, 1.0, 1.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0,

	-0.5, 0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0, 1.0, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 1.0, 0.0,
	-0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 0.0, 0.0,
	-0.5, 0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0,

	// black floor
	-1.0, -1.0, -0.5, 0.0, 0.0, 0.0, 0.0, 0.0,
	1.0, -1.0, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0,
	1.0, 1.0, -0.5, 0.0, 0.0, 0.0, 1.0, 1.0,
	1.0, 1.0, -0.5, 0.0, 0.0, 0.0, 1.0, 1.0,
	-1.0, 1.0, -0.5, 0.0, 0.0, 0.0, 0.0, 1.0,
	-1.0, -1.0, -0.5, 0.0, 0.0, 0.0, 0.0, 0.0,
}

var elements = []uint32{
	0, 1, 2,
	2, 3, 0,
}

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
