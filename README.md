# GoEngine3D

An early work in progress project on making a 3D graphics engine using [Go (Golang)](https://golang.org/) and [OpenGL](https://www.opengl.org/).

I have know knowledge on graphics programming, so this project is mainly a learning project for myself.

## Verisons / Packages
OpenGL v4.1

[go-gl/gl](https://github.com/go-gl/gl)

## Installation

### Windows
To get this application running on a Windows machine you will need to install a gcc compiler.
I used [MinGW](https://mingw-w64.org/doku.php), but there are many other options to use.
Also don't forget to link the gcc commands in your System Environment Variables.

### MacOS
MacOS is the main platform I have been developing on. I do not believe there is anything you needed to be installed before compiling.

### Linux
***Untested***


### Install Go (Golang)
Make sure to install [Go (Golang)](https://golang.org/) on your system. Instructions are found on its website.

### Compile commands
1. `go run *.go`
2. `go build -o <name>` then run your executable


# My Progress
### Creating a 3D plane - April 19th, 2019
![3D Plane](/documentation/3dPlane.png)

### Creating a reflective floor with stencils - April 23rd, 2019
![Reflective Floor](/documentation/reflectiveFloor.png)