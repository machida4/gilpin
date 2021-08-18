package gilpin

import (
	"bytes"
	"io"
)

type Encoder struct {
	writer io.Writer
	buffer *bytes.Buffer

	imageData *ImageData
}

func NewEncoder(writer io.Writer) *Encoder {
	buffer := new(bytes.Buffer)

	return &Encoder{writer: writer, buffer: buffer}
}

func (e *Encoder) Encode(imageData *ImageData) {
	e.imageData = imageData

	e.writeHeadData()
	e.writeCompressedData()
	e.writeTailData()
}

func (e *Encoder) writeHeadData() {
	e.writer.Write(e.imageData.headData)
}

func (e *Encoder) writeCompressedData() {
	e.writer.Write(e.imageData.compressedData)
}

func (e *Encoder) writeTailData() {
	e.writer.Write(e.imageData.tailData)
}
