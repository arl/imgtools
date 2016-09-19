package binimg

import (
	"image"
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
