package imstor_test

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/deiwin/imstor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	dataString = "somedata"
	data       = []byte(dataString)
	checksum   = "06343430109577305132"
	folderPath = "32/06343430109577305132"
	img        = image.NewGray16(image.Rect(0, 0, 3, 3))
	tempDir    string
	sizes      = []imstor.Size{
		imstor.Size{
			Name:   "small",
			Height: 30,
			Width:  30,
		}, imstor.Size{
			Name:   "large",
			Height: 300,
			Width:  300,
		},
	}
)

var _ = Describe("Imstor", func() {
	var s imstor.Storage
	BeforeEach(func() {
		var err error
		formats := []imstor.Format{png2JPEG{}, jpegFormat{}}
		tempDir, err = ioutil.TempDir("", "imstor-test")
		Expect(err).NotTo(HaveOccurred())
		conf := imstor.Config{
			RootPath:  tempDir,
			CopySizes: sizes,
		}
		s = imstor.NewWithCustomResizer(conf, formats, mockResizer{})
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Checksum", func() {
		It("should return the checksum for given bytes", func() {
			c := s.Checksum(data)
			Expect(c).To(Equal(checksum))
		})

		It("should be able to get the checksm for data encoded as a data URL", func() {
			c, err := s.ChecksumDataURL(fmt.Sprintf("data:,%s", dataString))
			Expect(err).NotTo(HaveOccurred())
			Expect(c).To(Equal(checksum))
		})
	})

	Describe("Store", func() {
		var expectImageFileToExist = func(name string) {
			path := filepath.Join(tempDir, filepath.FromSlash(folderPath), name)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				Fail(fmt.Sprintf("Expected file '%s' to exist", path))
			}
		}

		BeforeEach(func() {
			err := s.Store("image/jpeg", data)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create a image and copies", func() {
			expectImageFileToExist("original.jpg")
			expectImageFileToExist("small.jpg")
			expectImageFileToExist("large.jpg")
			// most assertions are in mock objects
		})

		It("should return proper path for the original image", func() {
			path, err := s.PathFor(checksum)
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(filepath.Join(filepath.FromSlash(folderPath), "original.jpg")))

		})

		It("should return proper paths for different sizes", func() {
			path, err := s.PathForSize(checksum, "small")
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(filepath.Join(filepath.FromSlash(folderPath), "small.jpg")))

			path, err = s.PathForSize(checksum, "large")
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(filepath.Join(filepath.FromSlash(folderPath), "large.jpg")))
		})

		It("should not return a path for improper size (say a prefix of an actual size)", func() {
			_, err := s.PathForSize(checksum, "smal")
			Expect(err).To(HaveOccurred())
		})
	})
})
