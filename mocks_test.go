package imstor_test

import (
	"image"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/deiwin/imstor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var smallImg = image.NewGray16(image.Rect(0, 0, 2, 2))
var largeImg = image.NewGray16(image.Rect(0, 0, 4, 4))
var newSizeImg = image.NewGray16(image.Rect(0, 0, 5, 5))

type png2JPEG struct {
	imstor.Format
}

func (f png2JPEG) DecodableMediaType() string {
	return "image/png"
}

type jpegFormat struct {
	imstor.Format
}

func (f jpegFormat) DecodableMediaType() string {
	return "image/jpeg"
}

func (f jpegFormat) Decode(r io.Reader) (image.Image, error) {
	bytes, err := ioutil.ReadAll(r)
	Expect(err).NotTo(HaveOccurred())
	Expect(bytes).To(Equal(data))
	return img, nil
}

func (f jpegFormat) EncodedExtension() string {
	return "jpg"
}

func (f jpegFormat) Encode(w io.Writer, i image.Image) error {
	if i == smallImg {
		expectToBeFile(w, "small.jpg")
	} else if i == largeImg {
		expectToBeFile(w, "large.jpg")
	} else if i == img {
		expectToBeFile(w, "original.jpg")
	} else if i == newSizeImg {
		expectToBeFile(w, "newFormat.jpg")
	} else {
		Fail("an unexpected image")
	}

	return nil
}

func expectToBeFile(w io.Writer, name string) {
	w.Write(data)
	path := filepath.Join(tempDir, filepath.FromSlash(folderPath), name)
	fileContents, err := ioutil.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())
	Expect(fileContents).To(Equal(data))
}

type mockResizer struct {
	imstor.Resizer
}

func (r mockResizer) Thumbnail(w, h uint, i image.Image) image.Image {
	Expect(i).To(Equal(img))
	if w == 30 && h == 30 {
		return smallImg
	} else if w == 300 && h == 300 {
		return largeImg
	} else if w == 16 && h == 16 {
		return newSizeImg
	}
	Fail("unexpected size")
	return nil
}
