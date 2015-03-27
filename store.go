package imstor

import (
	"bytes"
	"errors"
	"io"
	"log"

	"github.com/vincent-petithory/dataurl"
)

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
		if mediaType == format.EncodedExtension() {
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
