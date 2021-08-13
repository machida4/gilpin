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

func (e *Encoder) Encode()
