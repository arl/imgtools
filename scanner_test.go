package binimg

import (
	"image"
	"testing"
)

func testIsWhite(t *testing.T, newScanner func(image.Image) Scanner) {
	ss := []string{
		"000",
		"100",
		"011",
	}

	var testTbl = []struct {
		minx, miny, maxx, maxy int
		expected               bool
	}{
		{0, 0, 3, 3, false},
		{1, 1, 3, 3, false},
		{0, 1, 1, 2, true},
		{0, 0, 1, 1, false},
		{1, 0, 2, 1, false},
		{1, 0, 3, 2, false},
		{1, 2, 3, 3, true},
		{2, 2, 3, 3, true},
	}

	scanner := newScanner(newBinaryFromString(ss))
	for _, tt := range testTbl {
		actual := scanner.UniformColor(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy), White)
		if actual != tt.expected {
			t.Errorf("%d,%d|%d,%d): expected %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expected, actual)
		}
	}
}

func testIsBlack(t *testing.T, newScanner func(image.Image) Scanner) {
	ss := []string{
		"111",
		"011",
		"100",
	}

	var testTbl = []struct {
		minx, miny, maxx, maxy int
		expected               bool
	}{
		{0, 0, 3, 3, false},
		{1, 1, 3, 3, false},
		{0, 1, 1, 2, true},
		{0, 0, 1, 1, false},
		{1, 0, 2, 1, false},
		{1, 0, 3, 2, false},
		{1, 2, 3, 3, true},
		{2, 2, 3, 3, true},
	}

	scanner := newScanner(newBinaryFromString(ss))
	for _, tt := range testTbl {
		actual := scanner.UniformColor(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy), Black)
		if actual != tt.expected {
			t.Errorf("(%d,%d|%d,%d): expected %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expected, actual)
		}
	}
}

func testIsUniform(t *testing.T, newScanner func(image.Image) Scanner) {
	ss := []string{
		"111",
		"011",
		"100",
	}
	var testTbl = []struct {
		minx, miny, maxx, maxy int
		expected               bool
		expectedColor          Bit
	}{
		{0, 0, 3, 3, false, Bit{}},
		{1, 1, 3, 3, false, Bit{}},
		{0, 1, 1, 2, true, Black},
		{0, 0, 1, 1, true, White},
		{1, 0, 2, 1, true, White},
		{1, 0, 3, 2, true, White},
		{1, 2, 3, 3, true, Black},
		{2, 2, 3, 3, true, Black},
	}

	scanner := newScanner(newBinaryFromString(ss))
	for _, tt := range testTbl {
		actual, color := scanner.Uniform(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy))
		if actual != tt.expected {
			t.Errorf("(%d,%d|%d,%d): expected %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expected, actual)
		}
		if color != tt.expectedColor {
			t.Errorf("(%d,%d|%d,%d): expected color %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expectedColor, color)
		}
	}
}

func TestLinesScannerIsWhite(t *testing.T) {
	testIsWhite(t,
		func(img image.Image) Scanner {
			s, err := NewScanner(img)
			check(t, err)
			return s
		},
	)
}

func TestLinesScannerIsBlack(t *testing.T) {
	testIsBlack(t,
		func(img image.Image) Scanner {
			s, err := NewScanner(img)
			check(t, err)
			return s
		},
	)
}

func TestLinesScannerIsUniform(t *testing.T) {
	testIsUniform(t,
		func(img image.Image) Scanner {
			s, err := NewScanner(img)
			check(t, err)
			return s
		},
	)
}

func benchmarkScanner(b *testing.B, pngfile string, newScanner func(image.Image) Scanner) {
	img, err := loadPNG(pngfile)
	checkB(b, err)

	scanner := newScanner(NewFromImage(img))

	// run N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		scanner.UniformColor(img.Bounds(), White)
		scanner.UniformColor(img.Bounds(), Black)
		scanner.Uniform(img.Bounds())
	}
}

func BenchmarkLinesScanner(b *testing.B) {
	benchmarkScanner(b, "./testdata/big.png",
		func(img image.Image) Scanner {
			s, err := NewScanner(img)
			checkB(b, err)
			return s
		})
}
