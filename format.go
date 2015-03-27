package imstor

import (
	"image"
	"io"
)

type Format interface {
	DecodableMediaType() string
	Decode(io.Reader) (image.Image, error)
	Encode(io.Writer, image.Image) error
	EncodedExtension() string
}
