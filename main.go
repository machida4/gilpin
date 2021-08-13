package main

import (
	"fmt"
	"gilpin/decoder"
	"os"
	"path/filepath"
)

func main() {
	inputFilePath := filepath.Join("images", os.Args[1])
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	decoder := decoder.NewDecoder(inputFile)
	imageData := decoder.Decode()

	fmt.Println(imageData)

	fmt.Println("")
	fmt.Println("Complete")
}
