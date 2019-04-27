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
	viewDistance = 100.0
)

const SizeofFloat = 4.0

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	ship := Object{
		FileLocation: "./objects/shipWithNormals.obj",
	}
	vertices, elements := ship.Init()

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
	objectColor := gl.GetUniformLocation(program, gl.Str("objectColor\x00"))
	lightColor := gl.GetUniformLocation(program, gl.Str("lightColor\x00"))
	gl.Uniform3f(objectColor, 1, 0.5, 0.31)
	gl.Uniform3f(lightColor, 1, 1, 1)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 1.0, viewDistance)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{20, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 1})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// bind Frag Data Location
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Declare Vertex Array Object
	buildVAO()

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
	// texture, err := newTexture("./textures/diamonds.png")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// Position Attribute
	createPositionAttribute(program, "position\x00", 3, 3, 0)

	// Color Attribute
	// createColorAttribute(program, "color\x00", 3, 8, 3)

	// Texture Attribute
	// createTextureAttribute(program, "texcoord\x00", 2, 8, 6)

	// set clear window color
	gl.ClearColor(0, 0, 0, 1)

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
		time, elapsed := getTime(&previousTime)
		angle += elapsed

		updateWindowTitle(window, &time, &lastTime, &nbFrames)

		// render
		gl.UseProgram(program)

		// Draw cube and floor
		model = glm.HomogRotate3D(float32(angle), glm.Vec3{0.5, 0, 1})
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		// gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// gl.Enable(gl.STENCIL_TEST)

		// Draw Floor
		// gl.StencilFunc(gl.ALWAYS, 1, 0xFF)
		// gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE)
		// gl.StencilMask(0xFF)
		// gl.DepthMask(false)
		// gl.Clear(gl.STENCIL_BUFFER_BIT)

		// gl.DrawArrays(gl.TRIANGLES, 36, 6)

		// Draw reflection
		// gl.StencilFunc(gl.EQUAL, 1, 0xFF)
		// gl.StencilMask(0x00)
		// gl.DepthMask(true)

		// model = model.Mul4(glm.Translate3D(0, 0, -1))
		// gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		// gl.Uniform3f(uniColor, 0.3, 0.3, 0.3)
		// gl.DrawArrays(gl.TRIANGLES, 0, 36)
		// gl.Uniform3f(uniColor, 1, 1, 1)

		// gl.Disable(gl.STENCIL_TEST)

		// activate and bind texture
		// gl.ActiveTexture(gl.TEXTURE0)
		// gl.BindTexture(gl.TEXTURE_2D, texture)

		// gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))
		gl.DrawElements(gl.TRIANGLES, int32(len(elements)), gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func getTime(previousTime *float64) (float64, float64) {
	time := glfw.GetTime()
	elapsed := time - *previousTime
	*previousTime = time

	return time, elapsed
}

func updateWindowTitle(window *glfw.Window, time, lastTime *float64, nbFrames *int) {
	*nbFrames++
	if (*time - *lastTime) >= 1.0 {
		title := fmt.Sprintf("%v %v FPS | %v ms/frame", windowTitle, *nbFrames, 1000.0 / *nbFrames)
		window.SetTitle(title)
		*nbFrames = 0
		*lastTime += 1.0
	}
}
