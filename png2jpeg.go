package imstor

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

// PNG2JPEG format decodes an image from the PNG format and encodes it as a JPEG
var PNG2JPEG Format = png2JPEG{}

type png2JPEG struct {
}

func (f png2JPEG) Decode(r io.Reader) (image.Image, error) {
	return png.Decode(r)
}

func (f png2JPEG) DecodableMediaType() string {
	return "image/png"
}

func (f png2JPEG) Encode(w io.Writer, i image.Image) error {
	return jpeg.Encode(w, i, jpegEncodingOptions)
}

func (f png2JPEG) EncodedExtension() string {
	return "jpg"
}
