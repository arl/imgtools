package imgtools

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/aurelien-rainone/imgtools/internal/test"
)

func assertEqual(t *testing.T, a, b int) {
	if a != b {
		t.Fatalf("want %v == %v, got !=", a, b)
	}
}

func TestPowerOf2Roundup(t *testing.T) {
	assertEqual(t, pow2roundup(0), 1)
	assertEqual(t, pow2roundup(1), 1)
	assertEqual(t, pow2roundup(2), 2)
	assertEqual(t, pow2roundup(3), 4)
	assertEqual(t, pow2roundup(15), 16)
	assertEqual(t, pow2roundup(16), 16)
	assertEqual(t, pow2roundup(127), 128)
	assertEqual(t, pow2roundup(129), 256)
}

func TestPowerOf2Image(t *testing.T) {
	blue := color.RGBA{0, 0, 255, 255}
	red := color.RGBA{255, 0, 0, 255}
	uniform := image.NewUniform(blue)

	t.Run("pad a non power of 2 image", func(t *testing.T) {
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
			dst := PowerOf2Image(m, red)

			b := dst.Bounds()

			// width and height of new images should be equal
			assertEqual(t, b.Dx(), b.Dy())

			if pow2roundup(b.Dx()) != b.Dx() {
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
	t.Run("let power of 2 image as is", func(t *testing.T) {
		m := image.NewRGBA(image.Rect(0, 0, 16, 16))
		draw.Draw(m, m.Bounds(), uniform, image.ZP, draw.Src)
		dst := PowerOf2Image(m, red)

		b := dst.Bounds()

		// width and height of new images should be equal
		assertEqual(t, b.Dx(), b.Dy())

		if pow2roundup(b.Dx()) != b.Dx() {
			t.Errorf("dimension of new image should be a power of 2, got %v", b.Dx())
		}
		if m != dst.(*image.RGBA) {
			t.Errorf("want equal images, got different")
		}
		if err := test.Diff(m, dst); err != nil {
			t.Errorf("want same images bytes, got different: %v", err)
		}
	})
}
