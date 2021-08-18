package gilpin

import (
	"bytes"
	"io"
)

type Encoder struct {
	buffer *bytes.Buffer

	ImageData *ImageData
}

func NewEncoder(r io.Writer) *Encoder {

}

func (e *Encoder) Encode(imageData *ImageData) {

}

func (e *Encoder) writeHeadData() {}

func (e *Encoder) writeCompressedData() {}

func (e *Encoder) writeTailData() {}
