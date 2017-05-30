package binimg

import (
	"image"
	"image/color"
	"testing"

	"github.com/aurelien-rainone/imgtools/internal/test"
)

func TestIsOpaque(t *testing.T) {
	src, err := test.LoadPNG("../testdata/colorgopher.png")
	test.Check(t, err)

	bin := NewFromImage(src)
	if bin.Opaque() != true {
		t.Errorf("want Opaque to be true, got false")
	}
}

func TestSubImage(t *testing.T) {
	src, err := test.LoadPNG("../testdata/colorgopher.png")
	test.Check(t, err)

	sub := NewFromImage(src).SubImage(image.Rect(352, 352, 480, 480))
	refname := "../testdata/bwgopher.bottom-left.png"
	ref, err := test.LoadPNG(refname)
	test.Check(t, err)

	err = test.Diff(ref, sub)
	if err != nil {
		t.Errorf("converted image is different from %s: %v", refname, err)
	}
}

func TestEmptySubImage(t *testing.T) {
	empty := New(image.Rect(-5, -5, 10, 10)).SubImage(image.Rect(20, 20, 50, 50))
	if empty.Bounds() != New(image.Rectangle{}).Bounds() {
		t.Errorf("SubImage should produce an image with empty bounds when rects do not intersect")
	}
}

func TestPixelOperations(t *testing.T) {
	var (
		bin *Image
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))
	x, y := 9, 9

	blackRGBA := color.RGBA{0, 0, 0, 0xff}
	whiteRGBA := color.RGBA{0xff, 0xff, 0xff, 0xff}

	// get/set pixel from color.Color
	bin.Set(x, y, whiteRGBA)
	bit = Model.Convert(bin.At(x, y)).(Bit)
	if bit != On {
		t.Errorf("want bit at (%d,%d) to be On, got %v", x, y, bit)
	}

	// get/set pixel from color.Color
	bin.Set(x, y, blackRGBA)
	bit = Model.Convert(bin.At(x, y)).(Bit)
	if bit != Off {
		t.Errorf("want bit at (%d,%d) to be Off, got %v", x, y, bit)
	}

	// setting a pixel that is out of the image bounds should not panic, nor do nothing
	sub := bin.SubImage(image.Rect(1, 1, 2, 2)).(*Image)
	sub.Set(4, 4, whiteRGBA)
}

func TestBitOperations(t *testing.T) {
	var (
		bin *Image
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))
	x, y := 9, 9

	// get/set pixel from Bit
	bin.SetBit(x, y, On)
	bit = bin.BitAt(x, y)
	if bit != On {
		t.Errorf("expected pixel at (%d,%d) to be White, got %v", x, y, bit)
	}

	// get/set pixel from Bit
	bin.SetBit(x, y, Off)
	bit = bin.BitAt(x, y)
	if bit != Off {
		t.Errorf("expected pixel at (%d,%d) to be Black, got %v", x, y, bit)
	}

	// setting a bit that is out of the image bounds should not panic, nor do nothing
	sub := bin.SubImage(image.Rect(1, 1, 2, 2)).(*Image)
	sub.SetBit(4, 4, On)

	// getting a bit that is out of the image bound should return the zero
	// value of the color type
	bit = sub.BitAt(4, 4)
	var zero Bit
	if bit != zero {
		t.Errorf("expected BitAt to return Bit{} for out-of-bounds bit, got %v", bit)
	}
}

func TestColorModelIsComparable(t *testing.T) {
	bin := New(image.Rect(0, 0, 10, 10))
	if bin.ColorModel() != Model {
		t.Errorf("Binary.ColorModel should be comparable")
	}
}

func TestSetEmptyRect(t *testing.T) {
	var (
		bin *Image
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))

	// SetRect (empty rect)
	bin.SetRect(image.Rect(0, 0, 0, 0), On)

	// SetRect (rect which intersection is empty)
	bin.SetRect(image.Rect(100, 100, 10, 10), Off)
}

func TestSetRect(t *testing.T) {
	var (
		bin *Image
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))

	// SetRect
	bin.SetRect(image.Rect(0, 0, 1, 1), On)

	var testTbl = []struct {
		x, y int // x, y coordinates
		bit  Bit // expected color at specified coordinates
	}{
		{0, 0, On},
		{1, 0, Off},
		{0, 1, Off},
	}

	for _, tt := range testTbl {
		bit = bin.BitAt(tt.x, tt.y)
		if bit != tt.bit {
			t.Errorf("expected pixel at (%d,%d) to be %v, got %v", tt.x, tt.y, tt.bit, bit)
		}
	}
}

func TestSetOutOfBoundsRect(t *testing.T) {
	var (
		bin *Image
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))

	// setting a rect that goes out of the image bounds should not panic
	bin.SetRect(image.Rect(8, 8, 12, 12), On)

	var testTbl = []struct {
		x, y int // x, y coordinates
		bit  Bit // expected color at specified coordinates
	}{
		{8, 8, On},
		{9, 9, On},
	}

	for _, tt := range testTbl {
		bit = bin.BitAt(tt.x, tt.y)
		if bit != tt.bit {
			t.Errorf("expected pixel at (%d,%d) to be %v, got %v", tt.x, tt.y, tt.bit, bit)
		}
	}
}
