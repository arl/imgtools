package binimg

import (
	"bytes"
	"image"
)

// LinesScanner scans rectangular regions of a Binary image one line at at time.
type LinesScanner struct {
	bimg *Binary
}

// NewLinesScanner returns a new LinesScanner on the given Binary.
func NewLinesScanner(bimg *Binary) *LinesScanner {
	return &LinesScanner{bimg: bimg}
}

// UniformColor reports wether all the pixels of given region are of the color c.
func (s *LinesScanner) UniformColor(r image.Rectangle, c Bit) bool {
	// we want the other color for bytes.IndexBytes
	other := c.Other().v
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := s.bimg.PixOffset(r.Min.X, y)
		j := s.bimg.PixOffset(r.Max.X, y)
		// look for the first pixel that is not c
		if bytes.IndexByte(s.bimg.Pix[i:j], other) != -1 {
			return false
		}
	}
	return true
}

// Uniform reports wether the given region is uniform. If is the case, the
// uniform color bit is returned, otherwise the returned Bit is not
// significative (always the zero value of Bit).
func (s *LinesScanner) Uniform(r image.Rectangle) (bool, Bit) {
	// bit color of the first pixel (top-left)
	first := s.bimg.BitAt(r.Min.X, r.Min.Y)

	// check if all the pixels of the region are of this color.
	if s.UniformColor(r, first) {
		return true, first
	}
	return false, Bit{}
}
