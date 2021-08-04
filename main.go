package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gilpin/parser"
)

func main() {
	inputFilePath := filepath.Join("images", os.Args[1])
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	parser := parser.NewParser(inputFile)
	imageData := parser.Parse()

	fmt.Println(imageData)

	fmt.Println("")
	fmt.Println("Complete")
}
