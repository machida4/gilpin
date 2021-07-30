package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const pngSignature = "\x89PNG\r\n\x1a\n"

type Parser struct {
	buffer   *bytes.Buffer
	ihdr     *Ihdr
	seenIEND bool
}

type Ihdr struct {
	width, height int
	bitDepth      int
	colorType     int
	interlace     bool
}

func NewParser(r io.Reader) *Parser {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r)
	if err != nil {
		fmt.Println(err)
	}

	p := &Parser{buffer: buffer, ihdr: new(Ihdr), seenIEND: false}

	return p
}

func (p *Parser) next(n int) []byte {
	return p.buffer.Next(n)
}

func (p *Parser) parse() {
	p.checkSignature()

	for !p.seenIEND {
		p.parseChunk()
	}
}

func (p *Parser) checkSignature() {
	if string(p.next(8)) != pngSignature {
		fmt.Println("not PNG!!!")
	}
}

func (p *Parser) parseChunk() {
	length := int(binary.BigEndian.Uint32(p.next(4)))
	chunkType := string(p.next(4))

	switch chunkType {
	case "IHDR":
		fmt.Println("IHDR")
		p.parseIHDR(length)
		fmt.Println(p.ihdr)
	case "IPLT":
		fmt.Println("IPLT")
		p.readData(length)
	case "IDAT":
		fmt.Println("IDAT, ", length)
		p.readData(length)
	case "IEND":
		fmt.Println("IEND")
		p.seenIEND = true
	default:
		fmt.Println(chunkType)
		p.readData(length)
	}
}

func (p *Parser) parseIHDR(length int) {
	if length != 13 {
		fmt.Println("wrong IHDR format")
		return
	}

	p.ihdr.width = int(binary.BigEndian.Uint32(p.next(4)))
	p.ihdr.height = int(binary.BigEndian.Uint32(p.next(4)))
	p.ihdr.bitDepth = int(p.next(1)[0])
	p.ihdr.colorType = int(p.next(1)[0])

	compressionMethod := int(p.next(1)[0])
	if compressionMethod != 0 {
		fmt.Println("unknown compression method")
	}

	filterMethod := int(p.next(1)[0])
	if filterMethod != 0 {
		fmt.Println("unknown filter method")
	}

	p.ihdr.interlace = int(p.next(1)[0]) == 1

	p.readCRC()
}

func (p *Parser) readData(length int) {
	p.next(length)
	// crc
	p.readCRC()
}

func (p *Parser) readCRC() {
	p.next(4)
}

func main() {
	inputFilePath := filepath.Join("images", "lenna.png")
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	parser := NewParser(inputFile)
	parser.parse()

	fmt.Println("Complete")
}
