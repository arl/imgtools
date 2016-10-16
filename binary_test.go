package binimg

import (
	"image"
	"image/color"
	"testing"
)

func TestConvertImageThresholds(t *testing.T) {
	src, err := loadPNG("./testdata/colorgopher.png")
	check(t, err)
	var testTbl = []struct {
		m   binaryModel // binary model to use
		ref string      // reference file
	}{
		{BinaryModel, "./testdata/bwgopher.png"},
		{BinaryModelLowThreshold, "./testdata/bwgopher.low.threshold.png"},
		{BinaryModelHighThreshold, "./testdata/bwgopher.high.threshold.png"},
	}

	for _, tt := range testTbl {

		dst := NewCustomFromImage(src, tt.m)
		ref, _ := loadPNG(tt.ref)

		err = diff(ref, dst)
		if err != nil {
			t.Errorf("converted image is different from %s: %v", tt.ref, err)
		}
	}
}

func TestIsOpaque(t *testing.T) {
	src, err := loadPNG("./testdata/colorgopher.png")
	check(t, err)

	bin := NewFromImage(src)
	if bin.Opaque() != true {
		t.Errorf("expected Opaque to be true, got false")
	}
}

func TestSubImage(t *testing.T) {
	src, err := loadPNG("./testdata/colorgopher.png")
	check(t, err)

	sub := NewFromImage(src).SubImage(image.Rect(352, 352, 480, 480))
	refname := "./testdata/bwgopher.bottom-left.png"
	ref, _ := loadPNG(refname)

	err = diff(ref, sub)
	if err != nil {
		t.Errorf("converted image is different from %s: %v", refname, err)
	}
}

func TestBinaryModelConvert(t *testing.T) {
	var testTbl = []struct {
		m        binaryModel // binary model to use
		name     string      // model name (for test log)
		col      color.Color // source color component
		expected Bit         // expected converted Bit
	}{
		{BinaryModel, "BinaryModel", color.RGBA{0, 0, 0, 0}, Black},
		{BinaryModel, "BinaryModel", color.RGBA{0xff, 0xff, 0xff, 0xff}, White},
		// TODO: add others...
	}

	for _, tt := range testTbl {

		bit := tt.m.Convert(tt.col)
		if bit != tt.expected {
			t.Errorf("expected %s to convert %v to Bit value %v, got %v instead", tt.name, tt.col, tt.expected, bit)
		}
	}
}

func TestPixelOperations(t *testing.T) {
	var (
		bin *Binary
		bit Bit
		err error
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))
	x, y := 9, 9

	blackRGBA := color.RGBA{0, 0, 0, 0xff}
	whiteRGBA := color.RGBA{0xff, 0xff, 0xff, 0xff}

	// get/set pixel from color.Color
	bin.Set(x, y, whiteRGBA)
	bit = BinaryModel.Convert(bin.At(x, y)).(Bit)
	if bit != White {
		t.Errorf("expected pixel at (%d,%d) to be White, got %v", x, y, bit)
	}

	// get/set pixel from color.Color
	bin.Set(x, y, blackRGBA)
	bit = BinaryModel.Convert(bin.At(x, y)).(Bit)
	if bit != Black {
		t.Errorf("expected pixel at (%d,%d) to be Black, got %v", x, y, bit)
	}

	// setting a pixel that is out of the image bounds should not panic, nor do nothing
	sub := bin.SubImage(image.Rect(1, 1, 2, 2)).(*Binary)
	scanner, err := NewScanner(bin)
	check(t, err)

	sub.Set(4, 4, whiteRGBA)
	if !scanner.UniformColor(bin.Bounds(), Black) {
		t.Errorf("binary was expected to be uniformely black, got not uniform")
	}
}

func TestBitOperations(t *testing.T) {
	var (
		bin *Binary
		bit Bit
		err error
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))
	x, y := 9, 9

	// get/set pixel from Bit
	bin.SetBit(x, y, White)
	bit = bin.BitAt(x, y)
	if bit != White {
		t.Errorf("expected pixel at (%d,%d) to be White, got %v", x, y, bit)
	}

	// get/set pixel from Bit
	bin.SetBit(x, y, Black)
	bit = bin.BitAt(x, y)
	if bit != Black {
		t.Errorf("expected pixel at (%d,%d) to be Black, got %v", x, y, bit)
	}

	// setting a bit that is out of the image bounds should not panic, nor do nothing
	sub := bin.SubImage(image.Rect(1, 1, 2, 2)).(*Binary)
	scanner, err := NewScanner(bin)
	check(t, err)

	sub.SetBit(4, 4, White)
	if !scanner.UniformColor(bin.Bounds(), Black) {
		t.Errorf("binary was expected to be uniformely black, got not uniform")
	}
}

func TestSetEmptyRect(t *testing.T) {
	var (
		bin *Binary
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))

	// SetRect (empty rect)
	bin.SetRect(image.Rect(0, 0, 0, 0), White)

	// SetRect (rect which intersection is empty)
	bin.SetRect(image.Rect(100, 100, 10, 10), Black)
}

func TestSetRect(t *testing.T) {
	var (
		bin *Binary
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))

	// SetRect
	bin.SetRect(image.Rect(0, 0, 1, 1), White)

	var testTbl = []struct {
		x, y int // x, y coordinates
		bit  Bit // expected color at specified coordinates
	}{
		{0, 0, White},
		{1, 0, Black},
		{0, 1, Black},
	}

	for _, tt := range testTbl {
		bit = bin.BitAt(tt.x, tt.y)
		if bit != tt.bit {
			t.Errorf("expected pixel at (%d,%d) to be %s, got %v", tt.x, tt.y, tt.bit, bit)
		}
	}
}

func TestSetOutOfBoundsRect(t *testing.T) {
	var (
		bin *Binary
		bit Bit
	)

	// create a 10x10 Binary image
	bin = New(image.Rect(0, 0, 10, 10))

	// setting a rect that goes out of the image bounds should not panic
	bin.SetRect(image.Rect(8, 8, 12, 12), White)

	var testTbl = []struct {
		x, y int // x, y coordinates
		bit  Bit // expected color at specified coordinates
	}{
		{8, 8, White},
		{9, 9, White},
	}

	for _, tt := range testTbl {
		bit = bin.BitAt(tt.x, tt.y)
		if bit != tt.bit {
			t.Errorf("expected pixel at (%d,%d) to be %s, got %v", tt.x, tt.y, tt.bit, bit)
		}
	}
}
