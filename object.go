package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Object is a 3D model
type Object struct {
	FileLocation string
	Vertices     []float32
	Elements     []uint32
	Normals      []float32
}

// Init object
func (o Object) Init() ([]float32, []uint32) {
	vert, elem, norm := parseObjFile(o.FileLocation)

	var vertices []float32
	for i := 2; i < len(vert); i += 3 {
		vertices = append(vertices, vert[i-2])
		vertices = append(vertices, vert[i-1])
		vertices = append(vertices, vert[i])

		// vertices = append(vertices, norm[i-2])
		// vertices = append(vertices, norm[i-1])
		// vertices = append(vertices, norm[i])
	}

	o.Vertices = vertices
	o.Elements = elem
	fmt.Println("vertices", len(vert), vert)
	fmt.Println("elements", len(elem), elem)
	fmt.Println("normals", len(norm), norm)
	return vertices, elem
}

// Render is to be run during main loop
func (o Object) Render() {

}

func parseObjFile(filePath string) ([]float32, []uint32, []float32) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var vertices []float32
	var elements []uint32
	var normals []float32

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " ")
		if splitLine[0] == "v" {
			for i := 1; i < len(splitLine); i++ {
				val, _ := strconv.ParseFloat(splitLine[i], 4)
				vertices = append(vertices, float32(val))
			}
		} else if splitLine[0] == "f" {
			for i := 1; i < len(splitLine); i++ {
				splitObj := strings.Split(splitLine[i], "/")
				val, _ := strconv.Atoi(splitObj[0])
				elements = append(elements, uint32(val-1))

				// val, _ = strconv.Atoi(splitObj[2])
				// elements = append(elements, uint32(val-1))
			}
		} else if splitLine[0] == "vn" {
			for i := 1; i < len(splitLine); i++ {
				val, _ := strconv.ParseFloat(splitLine[i], 4)
				normals = append(normals, float32(val))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return vertices, elements, normals
}
