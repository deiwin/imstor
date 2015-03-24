// Package imstor enables you to create copies (or thumbnails) of your images and stores
// them along with the original image on your filesystem. The image and its
// copies are are stored in a file structure based on the (zero-prefixed, decimal)
// CRC 64 checksum of the original image. The last 2 characters of the checksum
// are used as the lvl 1 directory name.
//
// Example folder name and contents, given this checksum: 08446744073709551615:
// /configured/root/path/15/08446744073709551615/original.jpeg
// /configured/root/path/15/08446744073709551615/small.jpeg
// /configured/root/path/15/08446744073709551615/large.jpeg
package imstor

import (
	"hash/crc64"
	"image"
	"io"
)

var crcTable = crc64.MakeTable(crc64.ISO)

const (
	originalImageName = "original"
	pngExtension      = "png"
	jpegExtension     = "jpeg"
	permission        = 0640
)

type storage struct {
	conf    Config
	formats []Format
}

type Format interface {
	Encode(io.Writer, image.Image) error
	Decode(io.Reader) (image.Image, error)
	MediaType() string
	Extension() string
}

type Storage interface {
	StoreDataURLString(str string) error
	Store(mediaType string, data []byte) error
}

func New(conf Config, formats []Format) Storage {
	return storage{
		conf:    conf,
		formats: formats,
	}
}
