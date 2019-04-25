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
}

// Init object
func (o Object) Init() ([]float32, []uint32) {
	vert, elem := parseObjFile(o.FileLocation)
	fmt.Println(vert)
	fmt.Println(elem)
	return vert, elem
}

// Render is to be run during main loop
func (o Object) Render() {

}

func parseObjFile(filePath string) ([]float32, []uint32) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var vertices []float32
	var elements []uint32

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
				val, _ := strconv.Atoi(splitLine[i])
				elements = append(elements, uint32(val-1))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return vertices, elements
}
