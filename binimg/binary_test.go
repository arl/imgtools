package binimg

import (
	"image"
	"image/color"
	"testing"

	"github.com/aurelien-rainone/imgtools/internal/test"
)

func TestConvertImagePalette(t *testing.T) {
	src, err := test.LoadPNG("../testdata/colorgopher.png")
	test.Check(t, err)
	var testTbl = []struct {
		p   Palette // the palette to use
		ref string  // reference file
	}{
		{BlackAndWhite, "../testdata/bwgopher.png"},
		{BlackAndWhiteLowThreshold, "../testdata/bwgopher.low.threshold.png"},
		{BlackAndWhiteHighThreshold, "../testdata/bwgopher.high.threshold.png"},
		{Palette{97, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}}, "../testdata/redblue.gopher.png"},
	}

	for _, tt := range testTbl {
		dst := NewFromImage(src, tt.p)
		ref, err := test.LoadPNG(tt.ref)
		test.Check(t, err)

		err = test.Diff(ref, dst)
		if err != nil {
			t.Errorf("converted image is different from %s: %v", tt.ref, err)
		}
	}
}

func TestIsOpaque(t *testing.T) {
	src, err := test.LoadPNG("../testdata/colorgopher.png")
	test.Check(t, err)

	bin := NewFromImage(src, BlackAndWhite)
	if bin.Opaque() != true {
		t.Errorf("expected Opaque to be true, got false")
	}
}

func TestSubImage(t *testing.T) {
	src, err := test.LoadPNG("../testdata/colorgopher.png")
	test.Check(t, err)

	sub := NewFromImage(src, BlackAndWhite).SubImage(image.Rect(352, 352, 480, 480))
	refname := "../testdata/bwgopher.bottom-left.png"
	ref, err := test.LoadPNG(refname)
	test.Check(t, err)

	err = test.Diff(ref, sub)
	if err != nil {
		t.Errorf("converted image is different from %s: %v", refname, err)
	}
}

func TestPixelOperations(t *testing.T) {
	var (
		bin *Binary
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10), BlackAndWhite)
	x, y := 9, 9

	blackRGBA := color.RGBA{0, 0, 0, 0xff}
	whiteRGBA := color.RGBA{0xff, 0xff, 0xff, 0xff}

	// get/set pixel from color.Color
	bin.Set(x, y, whiteRGBA)
	bit = bin.Palette.ConvertBit(bin.At(x, y))
	if bit != On {
		t.Errorf("want bit at (%d,%d) to be On, got %v", x, y, bit)
	}

	// get/set pixel from color.Color
	bin.Set(x, y, blackRGBA)
	bit = bin.Palette.ConvertBit(bin.At(x, y))
	if bit != Off {
		t.Errorf("want bit at (%d,%d) to be Off, got %v", x, y, bit)
	}

	// setting a pixel that is out of the image bounds should not panic, nor do nothing
	sub := bin.SubImage(image.Rect(1, 1, 2, 2)).(*Binary)
	sub.Set(4, 4, whiteRGBA)
}

func TestBitOperations(t *testing.T) {
	var (
		bin *Binary
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10), BlackAndWhite)
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
	sub := bin.SubImage(image.Rect(1, 1, 2, 2)).(*Binary)
	sub.SetBit(4, 4, On)
}

func TestSetEmptyRect(t *testing.T) {
	var (
		bin *Binary
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10), BlackAndWhite)

	// SetRect (empty rect)
	bin.SetRect(image.Rect(0, 0, 0, 0), On)

	// SetRect (rect which intersection is empty)
	bin.SetRect(image.Rect(100, 100, 10, 10), Off)
}

func TestSetRect(t *testing.T) {
	var (
		bin *Binary
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10), BlackAndWhite)

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
		bin *Binary
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10), BlackAndWhite)

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
