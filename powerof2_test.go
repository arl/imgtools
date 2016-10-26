package imgtools

import (
	"image"
	"image/color"
	"testing"

	"github.com/aurelien-rainone/image/internal/test"
)

func assertEqual(t *testing.T, a, b int) {
	if a != b {
		t.Fatalf("expected %v == %v, got !=", a, b)
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
	var (
		src, dst image.Image
		err      error
	)
	src, err = test.loadPNG("./testdata/colorgopher.png")
	check(t, err)

	pad := color.RGBA{255, 0, 0, 255}
	dst = PowerOf2Image(src, pad)

	b := dst.Bounds()

	assertEqual(t, b.Dx(), b.Dy())

	corner := dst.At(b.Max.X-1, b.Max.Y-1)
	if corner != pad {
		t.Errorf("lower-left pixel should have the %v color, got %v instead", pad, corner)
	}
}

func TestAlreadyPowerOf2Image(t *testing.T) {
	var (
		src, dst image.Image
		err      error
	)
	src, err = test.loadPNG("./testdata/bwgopher.bottom-left.png")
	check(t, err)

	pad := color.RGBA{255, 0, 0, 255}
	dst = PowerOf2Image(src, pad)

	b := dst.Bounds()

	assertEqual(t, b.Dx(), b.Dy())
	if src != dst {
		t.Errorf("expecting PowerOf2Image to return src, as it is already a power-of-2 square")
	}
}
