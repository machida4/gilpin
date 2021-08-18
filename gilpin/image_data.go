package gilpin

import "fmt"

type ImageData struct {
	width, height int
	bitDepth      int
	colorType     int
	interlace     bool

	headData       []byte
	compressedData []byte
	tailData       []byte
}

type Scanline struct {
	filterType FilterType

	data []byte
}

func (id ImageData) String() string {
	return fmt.Sprintf("width: %d, height: %d, bitDepth: %d, colorType: %d, interlace: %t, headData: %d, compressedData: %d, tailData: %d",
		id.width, id.height, id.bitDepth, id.colorType, id.interlace, len(id.headData), len(id.compressedData), len(id.tailData),
	)
}

func (id *ImageData) CompressedData() []byte {
	return id.compressedData
}

func (id *ImageData) SetCompressedData(newCompressedData []byte) {
	id.compressedData = newCompressedData
}

func (id *ImageData) ToScanlines(filteredData []byte) []*Scanline {
	var scanlines []*Scanline
	scanlineSize := 1 + (bitPerPixel(id.colorType, id.bitDepth)*id.width+7)/8

	for h := 0; h < id.height; h++ {
		offset := h * scanlineSize
		filterType := FilterType(filteredData[offset])
		scanlineData := filteredData[offset+1 : offset+scanlineSize]

		scanline := &Scanline{filterType: filterType, data: scanlineData}
		scanlines = append(scanlines, scanline)
	}

	return scanlines
}
