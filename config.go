package imstor

import (
	"image"
	"io"

	"github.com/deiwin/gonfigure"
)

var (
	rootPathEnvProperty = gonfigure.NewRequiredEnvProperty("IMSTOR_ROOT_PATH")
)

type Config struct {
	RootPath  string
	CopySizes []Size
	Formats   []Format
}

func NewConfig(copySizes []Size, formats []Format) *Config {
	return &Config{
		RootPath:  rootPathEnvProperty.Value(),
		CopySizes: copySizes,
		Formats:   formats,
	}
}

// Size specifies a set of dimensions and a name that a copy of an image will
// be stored as
type Size struct {
	Name   string
	Height uint
	Width  uint
}

// A Format describes how an image of a certaing mimetype can be decoded and
// then encoded.
type Format interface {
	DecodableMediaType() string
	Decode(io.Reader) (image.Image, error)
	Encode(io.Writer, image.Image) error
	EncodedExtension() string
}
