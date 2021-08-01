package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const pngSignature = "\x89PNG\r\n\x1a\n"

type FilterType int

const (
	None = iota
	Sub
	Up
	Average
	Paeth
	Unknown
)

func (ft FilterType) String() string {
	return [...]string{"None", "Sub", "Up", "Average", "Paeth", "Unknown"}[ft]
}

type Parser struct {
	buffer   *bytes.Buffer
	seenIEND bool

	width, height int
	bitDepth      int
	colorType     int
	interlace     bool

	compressedData []byte
	scanlines      []*Scanline
}

type Scanline struct {
	filterType FilterType

	data []byte
}

func NewParser(r io.Reader) *Parser {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r)
	if err != nil {
		fmt.Println(err)
	}

	p := &Parser{buffer: buffer, seenIEND: false}

	return p
}

func (p *Parser) next(n int) []byte {
	return p.buffer.Next(n)
}

func (p *Parser) Parse() {
	p.checkSignature()

	for !p.seenIEND {
		p.parseChunk()
	}

	p.scanlines = p.divideFilteredDataIntoScanlines(inflate(p.compressedData))

	for _, scanline := range p.scanlines {
		fmt.Println("filter:", scanline.filterType)
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
		fmt.Printf("%+v\n", p)
	case "IPLT":
		fmt.Println("IPLT")
		p.skipData(length)
	case "IDAT":
		fmt.Println("IDAT,", length)
		p.parseIDAT(length)
	case "IEND":
		fmt.Println("IEND")
		p.seenIEND = true
	default:
		fmt.Println(chunkType)
		p.skipData(length)
	}
}

func (p *Parser) parseIHDR(length int) {
	if length != 13 {
		fmt.Println("wrong IHDR format")
		return
	}

	p.width = int(binary.BigEndian.Uint32(p.next(4)))
	p.height = int(binary.BigEndian.Uint32(p.next(4)))
	p.bitDepth = int(p.next(1)[0])
	p.colorType = int(p.next(1)[0])

	if int(p.next(1)[0]) != 0 {
		fmt.Println("unknown compression method")
	}
	if int(p.next(1)[0]) != 0 {
		fmt.Println("unknown filter method")
	}

	p.interlace = int(p.next(1)[0]) == 1

	p.skipCRC()
}

func (p *Parser) parseIDAT(length int) {
	p.compressedData = append(p.compressedData, p.next(length)...)
	p.skipCRC()
}

func (p *Parser) skipData(length int) {
	p.next(length)
	p.skipCRC()
}

func (p *Parser) skipCRC() {
	p.next(4)
}

func (p *Parser) divideFilteredDataIntoScanlines(filteredData []byte) []*Scanline {
	var scanlines []*Scanline
	scanlineSize := 1 + (bitPerPixel(p.colorType, p.bitDepth)*p.width+7)/8

	for h := 0; h < p.height; h++ {
		offset := h * scanlineSize
		filterType := FilterType(filteredData[offset])
		scanlineData := filteredData[offset+1 : offset+scanlineSize]

		scanline := &Scanline{filterType: filterType, data: scanlineData}
		scanlines = append(scanlines, scanline)
	}

	return scanlines
}

func bitPerPixel(colorType, depth int) int {
	switch colorType {
	case 0:
		return depth
	case 2:
		return depth * 3
	case 3:
		return depth
	case 4:
		return depth * 2
	case 6:
		return depth * 4
	default:
		return 0
	}
}

func inflate(data []byte) []byte {
	dataBuffer := bytes.NewReader(data)
	r, _ := zlib.NewReader(dataBuffer)
	defer r.Close()

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r)

	return buffer.Bytes()
}

func main() {
	inputFilePath := filepath.Join("images", os.Args[1])
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inputFile.Close()

	parser := NewParser(inputFile)
	parser.Parse()

	fmt.Println("Complete")
}
