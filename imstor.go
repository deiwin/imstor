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
	resizer Resizer
}

// Storage is the engine that can be used to store images and retrieve their paths
type Storage interface {
	StoreDataURLString(str string) error
	Store(mediaType string, data []byte) error
}

// New creates a storage engine using the default Resizer
func New(conf Config, formats []Format) Storage {
	return storage{
		conf:    conf,
		formats: formats,
	}
}

// NewWithCustomResizer creates a storage engine using a custom resizer
func NewWithCustomResizer(conf Config, formats []Format, resizer Resizer) Storage {
	return storage{
		conf:    conf,
		formats: formats,
		resizer: resizer,
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
