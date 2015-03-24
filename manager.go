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
	"fmt"
	"hash/crc64"
	"image"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/nfnt/resize"
	"github.com/vincent-petithory/dataurl"
)

// XXX
var conf Config
var formats []Format

var crcTable = crc64.MakeTable(crc64.ISO)

const (
	originalImageName = "original"
	pngExtension      = "png"
	jpegExtension     = "jpeg"
	permission        = 0640
)

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

func StoreDataURLString(s string) error {
	dataURL, err := dataurl.DecodeString(s)
	if err != nil {
		return err
	}
	return storeDataURL(dataURL)
}

func storeDataURL(dataURL *dataurl.DataURL) error {
	return Store(dataURL.MediaType.ContentType(), dataURL.Data)
}

func Store(mediaType string, data []byte) error {
	dataReader := bytes.NewReader(data)
	checksum := getChecksum(data)
	for _, format := range formats {
		if mediaType == format.MediaType() {
			return store(dataReader, checksum, format)
		}
	}
	return errors.New("Not a supported format!")
}

func store(r io.Reader, checksum string, f Format) error {
	image, err := f.Decode(r)
	if err != nil {
		return err
	}
	copies := createCopies(image, conf.CopySizes)
	folderPath := getAbsFolderPath(checksum)
	if err = createFolder(folderPath); err != nil {
		return err
	}
	if err = writeImageAndCopies(folderPath, image, copies, f); err != nil {
		log.Println("Writing an image failed, but a new folder and some files may have already been created. Please check your filesystem for clutter.")
		return err
	}
	return nil
}

func createCopies(image image.Image, sizes []Size) []imageFile {
	copies := make([]imageFile, len(sizes))
	for i, size := range sizes {
		imageCopy := resize.Thumbnail(size.Width, size.Height, image, resize.Lanczos3)
		copies[i] = imageFile{
			name:  size.Name,
			image: imageCopy,
		}
	}
	return copies
}

func writeImageAndCopies(folder string, original image.Image, copies []imageFile, f Format) error {
	imageFiles := append(copies, imageFile{
		name:  originalImageName,
		image: original,
	})
	return writeImageFiles(folder, imageFiles, f)
}

func writeImageFiles(folder string, imageFiles []imageFile, f Format) error {
	for _, imageFile := range imageFiles {
		fileName := fmt.Sprintf("%s.%s", imageFile.name, f.Extension())
		path := filepath.Join(folder, fileName)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, permission)
		if err != nil {
			return err
		}
		if err = f.Encode(file, imageFile.image); err != nil {
			return err
		}
	}
	return nil
}

func getAbsFolderPath(checksum string) string {
	structuredFolderPath := filepath.FromSlash(getStructuredFolderPath(checksum))
	return filepath.Join(conf.RootPath, structuredFolderPath)
}

func getStructuredFolderPath(checksum string) string {
	lvl1Dir := checksum[len(checksum)-2:]
	return path.Join(lvl1Dir, checksum)
}

func getChecksum(data []byte) string {
	crc := crc64.Checksum(data, crcTable)
	return fmt.Sprintf("%020d", crc)
}

func createFolder(path string) error {
	// rw-r-----
	return os.MkdirAll(path, permission)
}
