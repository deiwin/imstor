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
	formats = []imstor.Format{
		png2JPEG{},
		jpegFormat{},
	}
)

var _ = Describe("Imstor", func() {
	var s imstor.Storage
	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "imstor-test")
		Expect(err).NotTo(HaveOccurred())
		conf := &imstor.Config{
			RootPath:  tempDir,
			CopySizes: sizes,
			Formats:   formats,
		}
		s = imstor.NewWithCustomResizer(conf, mockResizer{})
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

		Context("with new configuration size added", func() {
			BeforeEach(func() {
				updatedSizes := append(sizes, imstor.Size{
					Name:   "newFormat",
					Height: 16,
					Width:  16,
				})
				conf := &imstor.Config{
					RootPath:  tempDir,
					CopySizes: updatedSizes,
					Formats:   formats,
				}
				s = imstor.NewWithCustomResizer(conf, mockResizer{})
			})

			Describe("storing the same image", func() {
				var err error
				BeforeEach(func() {
					err = s.Store("image/jpeg", data)
				})

				It("should return without an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should still have the image and copies plus the new one", func() {
					expectImageFileToExist("original.jpg")
					expectImageFileToExist("small.jpg")
					expectImageFileToExist("large.jpg")
					expectImageFileToExist("newFormat.jpg")
				})
			})
		})

		Describe("storing the same image", func() {
			var err error
			BeforeEach(func() {
				err = s.Store("image/jpeg", data)
			})

			It("should return without an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should still have the image and copies", func() {
				expectImageFileToExist("original.jpg")
				expectImageFileToExist("small.jpg")
				expectImageFileToExist("large.jpg")
			})
		})

		It("should return proper path for the original image", func() {
			path, err := s.PathFor(checksum)
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(filepath.Join(filepath.FromSlash(folderPath), "original.jpg")))
		})

		It("should return an error for an improper checksum", func() {
			_, err := s.PathFor("somethingrandom")
			Expect(err).To(HaveOccurred())
		})

		It("should return proper paths for different sizes", func() {
			path, err := s.PathForSize(checksum, "small")
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(filepath.Join(filepath.FromSlash(folderPath), "small.jpg")))

			path, err = s.PathForSize(checksum, "large")
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(filepath.Join(filepath.FromSlash(folderPath), "large.jpg")))
		})

		Describe("HasSizesForChecksum", func() {
			It("should return true if the sizes exist for that checksum", func() {
				hasSizes, err := s.HasSizesForChecksum(checksum, []string{"small", "large"})
				Expect(err).NotTo(HaveOccurred())
				Expect(hasSizes).To(BeTrue())
			})

			It("should return false if any of the sizes don't exist for that checksum", func() {
				hasSizes, err := s.HasSizesForChecksum(checksum, []string{"smallish", "large"})
				Expect(err).NotTo(HaveOccurred())
				Expect(hasSizes).To(BeFalse())
			})

			It("should return false for a random checksum", func() {
				hasSizes, err := s.HasSizesForChecksum("arandomchecksum", []string{"small", "large"})
				Expect(err).NotTo(HaveOccurred())
				Expect(hasSizes).To(BeFalse())
			})
		})

		It("should not return a path for improper size (say a prefix of an actual size)", func() {
			_, err := s.PathForSize(checksum, "smal")
			Expect(err).To(HaveOccurred())
		})
	})
})
