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
	"fmt"
	"hash/crc64"
	"path"
)

var crcTable = crc64.MakeTable(crc64.ISO)

const (
	originalImageName = "original"
)

type storage struct {
	conf    Config
	formats []Format
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

func getStructuredFolderPath(checksum string) string {
	lvl1Dir := checksum[len(checksum)-2:]
	return path.Join(lvl1Dir, checksum)
}

func getChecksum(data []byte) string {
	crc := crc64.Checksum(data, crcTable)
	return fmt.Sprintf("%020d", crc)
}
