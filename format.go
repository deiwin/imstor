package imstor

import (
	"image"
	"io"
)

// A Format describes how an image of a certaing mimetype can be decoded and
// then encoded.
type Format interface {
	DecodableMediaType() string
	Decode(io.Reader) (image.Image, error)
	Encode(io.Writer, image.Image) error
	EncodedExtension() string
}
