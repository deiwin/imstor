package imstor

import (
	"image"
	"image/jpeg"
	"io"
)

var jpegEncodingOptions = &jpeg.Options{
	Quality: jpeg.DefaultQuality,
}

var JPEGFormat Format = jpegFormat{}

type jpegFormat struct {
}

func (f jpegFormat) Decode(r io.Reader) (image.Image, error) {
	return jpeg.Decode(r)
}

func (f jpegFormat) DecodableMediaType() string {
	return "image/jpeg"
}

func (f jpegFormat) Encode(w io.Writer, i image.Image) error {
	return jpeg.Encode(w, i, jpegEncodingOptions)
}

func (f jpegFormat) EncodedExtension() string {
	return "jpg"
}
