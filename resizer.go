package imstor

import (
	"image"

	"github.com/nfnt/resize"
)

// A Resizer can resize an image into the given dimensions
type Resizer interface {
	// Resize should scale an image to new width and height. If one of the
	// parameters width or height is set to 0, its size will be calculated so that
	// the aspect ratio is that of the originating image.
	Resize(width, height uint, i image.Image) image.Image
	// Thumbnail should downscale provided image to max width and height preserving
	// original aspect ratio. It should return original image, without processing,
	// if original sizes are already smaller than the provided constraints.
	Thumbnail(maxWidth, maxHeight uint, i image.Image) image.Image
}

var DefaultResizer = defaultResizer{}

type defaultResizer struct{}

func (r defaultResizer) Resize(width, height uint, i image.Image) image.Image {
	return resize.Resize(width, height, i, resize.Lanczos3)
}

func (r defaultResizer) Thumbnail(maxWidth, maxHeight uint, i image.Image) image.Image {
	return resize.Thumbnail(maxWidth, maxHeight, i, resize.Lanczos3)
}
