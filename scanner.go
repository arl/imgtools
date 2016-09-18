package binimg

import (
	"image"
	"image/color"
)

// A Scanner is an object that can scan a rectangular colored region and report
// wether it is uniform (i.e composed of the same color) or not.
type Scanner interface {

	// UniformColor reports wether all the pixels of given region are of the color c.
	UniformColor(r image.Rectangle, c color.Color) bool

	// Uniform reports wether the given region is uniform. If is the case, the
	// uniform color is returned.
	Uniform(r image.Rectangle) (bool, color.Color)
}
