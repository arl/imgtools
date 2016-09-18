package binimg

import (
	"image"
	"image/color"
	"image/draw"
)

var (
	Black = Bit{0}
	White = Bit{255}
	Off   = Bit{0}
	On    = Bit{255}
)

// Bit represents a Black or White only binary color.
type Bit struct {
	v byte
}

func (c Bit) RGBA() (r, g, b, a uint32) {
	v := uint32(c.v)
	v |= v << 8
	return v, v, v, 0xffff
}

// Various binary models with different thresholds.
var (
	BinaryModelLowThreshold    binaryModel = NewBinaryModel(37)
	BinaryModelMediumThreshold binaryModel = NewBinaryModel(97)
	BinaryModelHighThreshold   binaryModel = NewBinaryModel(197)
	BinaryModel                binaryModel = BinaryModelMediumThreshold
)

type binaryModel struct {
	threshold uint8
}

func (m binaryModel) Convert(c color.Color) color.Color {
	if _, ok := c.(Bit); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	y := (299*r + 587*g + 114*b + 500) / 1000
	if uint8(y>>8) > m.threshold {
		return White
	}
	return Black
}

func NewBinaryModel(threshold uint8) binaryModel {
	return binaryModel{threshold}
}

// Binary is an in-memory image whose At method returns Bit values.
type Binary struct {
	// Pix holds the image's pixels, as 0 or 1 uint8 values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle

	model binaryModel
}

func (b *Binary) ColorModel() color.Model { return b.model }

func (b *Binary) Bounds() image.Rectangle { return b.Rect }

func (b *Binary) At(x, y int) color.Color {
	return b.BinaryAt(x, y)
}

func (b *Binary) BinaryAt(x, y int) Bit {
	if !(image.Point{x, y}.In(b.Rect)) {
		return Bit{}
	}
	i := b.PixOffset(x, y)
	return Bit{b.Pix[i]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (b *Binary) PixOffset(x, y int) int {
	return (y-b.Rect.Min.Y)*b.Stride + (x-b.Rect.Min.X)*1
}

func (b *Binary) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(b.Rect)) {
		return
	}
	i := b.PixOffset(x, y)
	b.Pix[i] = b.model.Convert(c).(Bit).v
}

func (b *Binary) SetBit(x, y int, c Bit) {
	if !(image.Point{x, y}.In(b.Rect)) {
		return
	}
	i := b.PixOffset(x, y)
	b.Pix[i] = c.v
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (b *Binary) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(b.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Binary{}
	}
	i := b.PixOffset(r.Min.X, r.Min.Y)
	return &Binary{
		Pix:    b.Pix[i:],
		Stride: b.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (b *Binary) Opaque() bool {
	return true
}

// NewBinary returns a new Binary image with the given bounds.
func NewBinary(r image.Rectangle) *Binary {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 1*w*h)
	return &Binary{pix, 1 * w, r, BinaryModel}
}

// NewCustomBinary returns a new Binary image with the given bounds and binary
// model.
func NewCustomBinary(r image.Rectangle, model binaryModel) *Binary {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 1*w*h)
	return &Binary{pix, 1 * w, r, model}
}

// NewFromImage returns the binary image that is the conversion of the given
// source image.
func NewFromImage(src image.Image) *Binary {
	dst := NewBinary(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, image.Point{}, draw.Src)
	return dst
}

// NewCustomFromImage returns the binary image that is the conversion of the
// given source image with the specified binary model.
func NewCustomFromImage(src image.Image, model binaryModel) *Binary {
	dst := NewCustomBinary(src.Bounds(), model)
	draw.Draw(dst, dst.Bounds(), src, image.Point{}, draw.Src)
	return dst
}
