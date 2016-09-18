package binimg

import "image"

// A Scanner is an object that can scan a rectangular colored region and report
// wether it is uniform (i.e composed of the same color) or not.
type Scanner interface {

	// UniformColor reports wether all the pixels of given region are of the color c.
	UniformColor(r image.Rectangle, c Bit) bool

	// Uniform reports wether the given region is uniform. If is the case, the
	// uniform color bit is returned, otherwise the returned Bit is not
	// significative (always the zero value of Bit).
	Uniform(r image.Rectangle) (bool, Bit)
}
