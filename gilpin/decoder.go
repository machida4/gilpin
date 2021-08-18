package gilpin

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

const pngSignature = "\x89PNG\r\n\x1a\n"

type Decoder struct {
	buffer   *bytes.Buffer
	seenIEND bool

	imageData *ImageData
}

func NewDecoder(r io.Reader) *Decoder {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r)
	if err != nil {
		fmt.Println(err)
	}

	return &Decoder{buffer: buffer, imageData: new(ImageData), seenIEND: false}
}

func (d *Decoder) Decode() *ImageData {
	d.checkSignature()

	for !d.seenIEND {
		d.parseChunk()
	}

	return d.imageData
}

func (d *Decoder) next(n int) []byte {
	return d.buffer.Next(n)
}

func (d *Decoder) checkSignature() {
	if string(d.next(8)) != pngSignature {
		fmt.Println("not PNG!!!")
	}
}

func (d *Decoder) parseChunk() {
	length := int(binary.BigEndian.Uint32(d.next(4)))
	chunkType := string(d.next(4))

	switch chunkType {
	case "IHDR":
		d.parseIHDR(length)
	case "IDAT":
		d.parseIDAT(length)
	case "IEND":
		d.seenIEND = true
	default:
		d.parseDefault(length)
	}
}

func (d *Decoder) parseIHDR(length int) {
	if length != 13 {
		fmt.Println("wrong IHDR format")
		return
	}

	IHDRData := d.next(13)
	d.imageData.headData = append(d.imageData.headData, IHDRData...)

	d.imageData.width = int(binary.BigEndian.Uint32(IHDRData[0:4]))
	d.imageData.height = int(binary.BigEndian.Uint32(IHDRData[4:8]))
	d.imageData.bitDepth = int(IHDRData[8])
	d.imageData.colorType = int(IHDRData[9])

	if int(IHDRData[10]) != 0 {
		fmt.Println("unknown compression method")
	}
	if int(IHDRData[11]) != 0 {
		fmt.Println("unknown filter method")
	}

	d.imageData.interlace = int(IHDRData[12]) == 1

	d.skipCRC()
}

func (d *Decoder) parseIDAT(length int) {
	d.imageData.compressedData = append(d.imageData.compressedData, d.next(length)...)
	d.skipCRC()
}

func (d *Decoder) parseDefault(length int) {
	if len(d.imageData.compressedData) == 0 {
		d.imageData.headData = append(d.imageData.headData, d.next(length)...)
	} else {
		d.imageData.tailData = append(d.imageData.tailData, d.next(length)...)
	}

	d.skipCRC()
}

func (d *Decoder) skipCRC() {
	d.next(4)
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
