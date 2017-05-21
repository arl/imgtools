package imgscan

import (
	"image"
	"image/color"
)

// A Scanner is an image that can report wether a rectangular region is uniform
// (i.e composed of the same color) or not.
type Scanner interface {
	image.Image

	// UniformColor reports wether all the pixels of given region are of the color c.
	UniformColor(r image.Rectangle, c color.Color) bool

	// Uniform reports wether the given region is uniform. If that is the case, the
	// uniform color is returned, otherwise the returned is nil.
	Uniform(r image.Rectangle) (bool, color.Color)
}
