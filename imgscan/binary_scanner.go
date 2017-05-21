package imgscan

import (
	"bytes"
	"fmt"
	"image"
	"image/color"

	"github.com/aurelien-rainone/imgtools/binimg"
)

// NewScanner returns a new Scanner of the given image.Image.
//
// The actual scanner implementation depends on the image bit depth and the
// availability of an implementation.
func NewScanner(img image.Image) (Scanner, error) {
	switch impl := img.(type) {
	case *binimg.Binary:
		return &binaryScanner{impl}, nil
	case *image.Alpha:
	case *image.Gray:
	default:
	}
	return nil, fmt.Errorf("unsupported image type")
}

type binaryScanner struct {
	*binimg.Binary
}

// UniformColor reports wether all the pixels of given region are of the color c.
//
// If c is not a color ot type binimg.Bit, UniformColor will panic.
func (s *binaryScanner) UniformColor(r image.Rectangle, c color.Color) bool {
	// in a binary image, pixel/bytes are 1 or 0, we want the other color for
	// bytes.IndexBytes
	other := c.(binimg.Bit).Other().V
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := s.PixOffset(r.Min.X, y)
		j := s.PixOffset(r.Max.X, y)
		// look for the first pixel that is not c
		if bytes.IndexByte(s.Pix[i:j], other) != -1 {
			return false
		}
	}
	return true
}

// Uniform reports wether the given region is uniform. If that is the case, the
// uniform color is returned, otherwise the returned color is nil.
func (s *binaryScanner) Uniform(r image.Rectangle) (bool, color.Color) {
	// bit color of the first pixel (top-left)
	first := s.BitAt(r.Min.X, r.Min.Y)

	// check if all the pixels of the region are of this color.
	if s.UniformColor(r, first) {
		return true, first
	}
	return false, nil
}
