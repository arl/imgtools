package binimg

import (
	"bytes"
	"image"
)

// LinesScanner scans rectangular regions of a Binary image one line at at time.
type LinesScanner struct {
	bimg *Binary
}

func NewLinesScanner(bimg *Binary) *LinesScanner {
	return &LinesScanner{bimg: bimg}
}

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

func (s *LinesScanner) Uniform(r image.Rectangle) (bool, Bit) {
	// bit color of the first pixel (top-left)
	first := s.bimg.BitAt(r.Min.X, r.Min.Y)

	// check if all the pixels of the region are of this color.
	if s.UniformColor(r, first) {
		return true, first
	}
	return false, Bit{}
}
