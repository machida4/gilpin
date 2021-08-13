package parser

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
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

	imageData *ImageData
}

type ImageData struct {
	width, height int
	bitDepth      int
	colorType     int
	interlace     bool

	headData       []byte
	compressedData []byte
	tailData       []byte

	// scanlines []*Scanline
}

func (id ImageData) String() string {
	return fmt.Sprintf("width: %d, height: %d, bitDepth: %d, colorType: %d, interlace: %t, headData: %d, compressedData: %d, tailData: %d",
		id.width, id.height, id.bitDepth, id.colorType, id.interlace, len(id.headData), len(id.compressedData), len(id.tailData),
	)
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

	p := &Parser{buffer: buffer, imageData: new(ImageData), seenIEND: false}

	return p
}

func (p *Parser) next(n int) []byte {
	return p.buffer.Next(n)
}

func (p *Parser) Parse() *ImageData {
	p.checkSignature()

	for !p.seenIEND {
		p.parseChunk()
	}

	// p.scanlines = p.divideFilteredDataIntoScanlines(inflate(p.compressedData))

	// for _, scanline := range p.scanlines {
	// 	fmt.Println("filter:", scanline.filterType)
	// }

	return p.imageData
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
		p.parseIHDR(length)
	case "IDAT":
		p.parseIDAT(length)
	case "IEND":
		p.seenIEND = true
	default:
		p.parseDefault(length)
	}
}

func (p *Parser) parseIHDR(length int) {
	if length != 13 {
		fmt.Println("wrong IHDR format")
		return
	}

	IHDRData := p.next(13)
	p.imageData.headData = append(p.imageData.headData, IHDRData...)

	p.imageData.width = int(binary.BigEndian.Uint32(IHDRData[0:4]))
	p.imageData.height = int(binary.BigEndian.Uint32(IHDRData[4:8]))
	p.imageData.bitDepth = int(IHDRData[8])
	p.imageData.colorType = int(IHDRData[9])

	if int(IHDRData[10]) != 0 {
		fmt.Println("unknown compression method")
	}
	if int(IHDRData[11]) != 0 {
		fmt.Println("unknown filter method")
	}

	p.imageData.interlace = int(IHDRData[12]) == 1

	p.skipCRC()
}

func (p *Parser) parseIDAT(length int) {
	p.imageData.compressedData = append(p.imageData.compressedData, p.next(length)...)
	p.skipCRC()
}

func (p *Parser) parseDefault(length int) {
	if len(p.imageData.compressedData) == 0 {
		p.imageData.headData = append(p.imageData.headData, p.next(length)...)
	} else {
		p.imageData.tailData = append(p.imageData.tailData, p.next(length)...)
	}

	p.skipCRC()
}

func (p *Parser) skipCRC() {
	p.next(4)
}

// func (p *Parser) divideFilteredDataIntoScanlines(filteredData []byte) []*Scanline {
// 	var scanlines []*Scanline
// 	scanlineSize := 1 + (bitPerPixel(p.imageData.colorType, p.imageData.bitDepth)*p.imageData.width+7)/8

// 	for h := 0; h < p.imageData.height; h++ {
// 		offset := h * scanlineSize
// 		filterType := FilterType(filteredData[offset])
// 		scanlineData := filteredData[offset+1 : offset+scanlineSize]

// 		scanline := &Scanline{filterType: filterType, data: scanlineData}
// 		scanlines = append(scanlines, scanline)
// 	}

// 	return scanlines
// }

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
