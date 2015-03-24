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
	"bytes"
	"errors"
	"hash/crc64"
	"image"
	"io"
	"log"

	"github.com/vincent-petithory/dataurl"
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

type imageFile struct {
	name  string
	image image.Image
}

type Format interface {
	Encode(io.Writer, image.Image) error
	Decode(io.Reader) (image.Image, error)
	MediaType() string
	Extension() string
}

func (s storage) StoreDataURLString(str string) error {
	dataURL, err := dataurl.DecodeString(str)
	if err != nil {
		return err
	}
	return s.storeDataURL(dataURL)
}

func (s storage) storeDataURL(dataURL *dataurl.DataURL) error {
	return s.Store(dataURL.MediaType.ContentType(), dataURL.Data)
}

func (s storage) Store(mediaType string, data []byte) error {
	dataReader := bytes.NewReader(data)
	checksum := getChecksum(data)
	for _, format := range s.formats {
		if mediaType == format.MediaType() {
			return s.storeInFormat(dataReader, checksum, format)
		}
	}
	return errors.New("Not a supported format!")
}

func (s storage) storeInFormat(r io.Reader, checksum string, f Format) error {
	image, err := f.Decode(r)
	if err != nil {
		return err
	}
	copies := createCopies(image, s.conf.CopySizes)
	folderPath := getAbsFolderPath(s.conf.RootPath, checksum)
	if err = createFolder(folderPath); err != nil {
		return err
	}
	if err = writeImageAndCopies(folderPath, image, copies, f); err != nil {
		log.Println("Writing an image failed, but a new folder and some files may have already been created. Please check your filesystem for clutter.")
		return err
	}
	return nil
}
