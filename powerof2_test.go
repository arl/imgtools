package imgtools

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"github.com/arl/imgtools/binimg"
	"github.com/arl/imgtools/internal/test"
)

func assertEqual(t *testing.T, a, b int) {
	if a != b {
		t.Fatalf("want %v == %v, got !=", a, b)
	}
}

func TestPowerOf2Roundup(t *testing.T) {
	assertEqual(t, Pow2Roundup(0), 1)
	assertEqual(t, Pow2Roundup(1), 1)
	assertEqual(t, Pow2Roundup(2), 2)
	assertEqual(t, Pow2Roundup(3), 4)
	assertEqual(t, Pow2Roundup(15), 16)
	assertEqual(t, Pow2Roundup(16), 16)
	assertEqual(t, Pow2Roundup(127), 128)
	assertEqual(t, Pow2Roundup(129), 256)
}

func TestPowerOf2Image(t *testing.T) {
	blue := color.RGBA{0, 0, 255, 255}
	red := color.RGBA{255, 0, 0, 255}
	uniform := image.NewUniform(blue)

	t.Run("new square image is padded", func(t *testing.T) {
		var tests = []struct {
			w, h int // org image dimensions
		}{
			{2, 3},
			{12, 14},
			{15, 17},
			{16, 1},
		}

		for _, tt := range tests {
			m := image.NewRGBA(image.Rect(0, 0, tt.w, tt.h))
			draw.Draw(m, m.Bounds(), uniform, image.ZP, draw.Src)
			dst, err := PowerOf2Image(m, red)
			test.Check(t, err)

			b := dst.Bounds()

			// width and height of new images should be equal
			assertEqual(t, b.Dx(), b.Dy())

			if Pow2Roundup(b.Dx()) != b.Dx() {
				t.Errorf("dimension of new image should be a power of 2, got %v", b.Dx())
			}

			if topLeft := dst.At(0, 0); topLeft != blue {
				t.Errorf("want top-left pixel color unchanged (blue), got %v", topLeft)
			}
			if bottomRight := dst.At(b.Dx()-1, b.Dy()-1); bottomRight != red {
				t.Errorf("want bottom-right pixel red (padding), got %v", bottomRight)
			}
		}
	})
	t.Run("do not touch image if already power of 2", func(t *testing.T) {
		m := image.NewRGBA(image.Rect(0, 0, 16, 16))
		draw.Draw(m, m.Bounds(), uniform, image.ZP, draw.Src)
		dst, err := PowerOf2Image(m, red)
		test.Check(t, err)

		b := dst.Bounds()

		// width and height of new images should be equal
		assertEqual(t, b.Dx(), b.Dy())

		if Pow2Roundup(b.Dx()) != b.Dx() {
			t.Errorf("dimension of new image should be a power of 2, got %v", b.Dx())
		}
		if m != dst.(*image.RGBA) {
			t.Errorf("want equal images, got different")
		}
		if err := test.Diff(m, dst); err != nil {
			t.Errorf("want same images bytes, got different: %v", err)
		}
	})
	t.Run("return a square image of the same type", func(t *testing.T) {
		var tests = []struct {
			org draw.Image
		}{
			{image.NewGray(image.Rect(2, 3, 4, 5))},
			{image.NewGray16(image.Rect(2, 3, 4, 5))},
			{binimg.New(image.Rect(0, -1, 12, 14))},
			{image.NewAlpha(image.Rect(2, 3, 4, 5))},
		}

		for _, tt := range tests {
			dst, err := PowerOf2Image(tt.org, red)
			test.Check(t, err)

			torg, tdst := reflect.TypeOf(tt.org), reflect.TypeOf(dst)
			if tdst != torg {
				t.Errorf("want new image of same type (%v), got (%v)", torg.Name(), tdst.Name())
			}
		}
	})
}
