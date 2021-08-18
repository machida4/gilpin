package main

import (
	"fmt"
	"gilpin/gilpin"
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

	decoder := gilpin.NewDecoder(inputFile)
	imageData := decoder.Decode()

	outputFilePath := filepath.Join("images", os.Args[2])
	outputFile, err := os.Open(outputFilePath)
	encoder := gilpin.NewEncoder(outputFile)
	encoder.Encode(imageData)

	fmt.Println(imageData)

	fmt.Println("")
	fmt.Println("Complete")
}
