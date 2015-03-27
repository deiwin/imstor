package imstor

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

const permission = 0640

type imageFile struct {
	name  string
	image image.Image
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
		fileName := fmt.Sprintf("%s.%s", imageFile.name, f.DecodableMediaType())
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

func getAbsFolderPath(rootPath string, checksum string) string {
	structuredFolderPath := filepath.FromSlash(getStructuredFolderPath(checksum))
	return filepath.Join(rootPath, structuredFolderPath)
}

func createFolder(path string) error {
	// rw-r-----
	return os.MkdirAll(path, permission)
}