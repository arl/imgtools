package imgscan

import (
	"image"
	"image/color"
	"testing"

	"github.com/aurelien-rainone/imgtools/binimg"
	"github.com/aurelien-rainone/imgtools/internal/test"
)

func newBinaryFromString(ss []string) *binimg.Binary {
	w, h := len(ss[0]), len(ss)
	for i := range ss {
		if len(ss[i]) != w {
			panic("all strings should have the same length")
		}
	}

	img := binimg.New(image.Rect(0, 0, w, h))
	for y := range ss {
		for x := range ss[y] {
			if ss[y][x] == '1' {
				img.SetBit(x, y, binimg.On)
			}
		}
	}
	return img
}

func TestBinaryScannerIsUniformColor(t *testing.T) {
	ss := []string{
		"000",
		"100",
		"011",
	}

	var tests = []struct {
		minx, miny, maxx, maxy int
		col                    color.Color
		uniform                bool
	}{
		{0, 0, 3, 3, binimg.On, false},
		{0, 0, 3, 3, binimg.Off, false},
		{1, 1, 3, 3, binimg.On, false},
		{1, 1, 3, 3, binimg.Off, false},
		{0, 1, 1, 2, binimg.On, true},
		{0, 1, 1, 2, binimg.Off, false},
		{0, 0, 1, 1, binimg.On, false},
		{0, 0, 1, 1, binimg.Off, true},
		{1, 0, 2, 1, binimg.On, false},
		{1, 0, 2, 1, binimg.Off, true},
		{1, 0, 3, 2, binimg.On, false},
		{1, 0, 3, 2, binimg.Off, true},
		{1, 2, 3, 3, binimg.On, true},
		{1, 2, 3, 3, binimg.Off, false},
		{2, 2, 3, 3, binimg.On, true},
		{2, 2, 3, 3, binimg.Off, false},
	}

	scanner, err := NewScanner(newBinaryFromString(ss))
	test.Check(t, err)
	for _, tt := range tests {
		got := scanner.IsUniformColor(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy), tt.col)
		if got != tt.uniform {
			t.Errorf("want %v for IsUniformColor(rect{%d,%d|%d,%d}, col:%v), got %v", tt.uniform, tt.minx, tt.miny, tt.maxx, tt.maxy, tt.col, got)
		}
	}
}

func TestBinaryScannerIsUniform(t *testing.T) {
	ss := []string{
		"000",
		"100",
		"011",
	}

	var tests = []struct {
		minx, miny, maxx, maxy int
		col                    color.Color
		uniform                bool
	}{
		{0, 0, 3, 3, nil, false},
		{1, 1, 3, 3, nil, false},
		{0, 1, 1, 2, binimg.On, true},
		{0, 0, 1, 1, binimg.Off, true},
		{1, 0, 2, 1, binimg.Off, true},
		{1, 0, 3, 2, binimg.Off, true},
		{1, 1, 2, 3, nil, false},
		{1, 2, 3, 3, binimg.On, true},
		{2, 2, 3, 3, binimg.On, true},
	}

	scanner, err := NewScanner(newBinaryFromString(ss))
	test.Check(t, err)
	for _, tt := range tests {
		got, col := scanner.IsUniform(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy))
		if got != tt.uniform {
			t.Errorf("want uniform=%v for IsUniform(rect{%d,%d|%d,%d}), got %v", tt.uniform, tt.minx, tt.miny, tt.maxx, tt.maxy, tt.col, got)
		}
		if col != tt.col {
			t.Errorf("want color=%v for IsUniform(rect{%d,%d|%d,%d}), got %v", tt.col, tt.minx, tt.miny, tt.maxx, tt.maxy, col)
		}
	}
}

func benchmarkScanner(b *testing.B, pngfile string, newScanner func(image.Image) Scanner) {
	img, err := test.LoadPNG(pngfile)
	test.CheckB(b, err)

	scanner := newScanner(binimg.NewFromImage(img))

	// run N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		scanner.IsUniformColor(img.Bounds(), color.White)
		scanner.IsUniformColor(img.Bounds(), color.Black)
		scanner.IsUniform(img.Bounds())
	}
}

func BenchmarkLinesScanner(b *testing.B) {
	benchmarkScanner(b, "./testdata/big.png",
		func(img image.Image) Scanner {
			s, err := NewScanner(img)
			test.CheckB(b, err)
			return s
		})
}
