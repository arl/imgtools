package imgscan

import (
	"encoding/csv"
	"image"
	"image/color"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/arl/imgtools/binimg"
	"github.com/arl/imgtools/internal/test"
)

func newGrayFromString(ss []string) *image.Gray {
	w, h := strings.Count(ss[0], ",")+1, len(ss)
	for i := range ss {
		if strings.Count(ss[i], ",")+1 != w {
			panic("all strings should have the same number of ','")
		}
	}

	var (
		err    error
		record []string
	)
	img := image.NewGray(image.Rect(0, 0, w, h))
	for y, s := range ss {
		r := csv.NewReader(strings.NewReader(s))
		r.TrimLeadingSpace = true
		record, err = r.Read()
		if err != nil {
			log.Fatalln("can't parse csv:", err)
		}
		for x, val := range record {
			if i, err := strconv.ParseInt(val, 0, 16); err != nil {
				log.Fatalf("Can't get value of pixel %v,%v: %v\n", x, y, err)
			} else {
				img.SetGray(x, y, color.Gray{uint8(i)})
			}
		}
	}
	return img
}

func TestGrayScannerIsUniformColor(t *testing.T) {
	ss := []string{
		"  0,   0,   0",
		"122,   0,   0",
		"  0,  24,  24",
	}

	var tests = []struct {
		minx, miny, maxx, maxy int
		col                    color.Color
		uniform                bool
	}{
		{0, 0, 3, 3, binimg.On, false},
		{0, 0, 3, 3, color.White, false},
		{0, 0, 3, 3, color.RGBA{12, 23, 34, 45}, false},
		{0, 1, 1, 2, color.Gray{122}, true},
		{1, 2, 3, 3, color.Gray{24}, true},
		{1, 2, 3, 3, color.Gray{127}, false},
	}

	img := newGrayFromString(ss)
	scanner, err := NewScanner(img)
	test.Check(t, err)
	for _, tt := range tests {
		uniform := scanner.IsUniformColor(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy), tt.col)
		if uniform != tt.uniform {
			t.Errorf("want %v for IsUniformColor(rect{%d,%d|%d,%d}, col:%v), got %v", tt.uniform, tt.minx, tt.miny, tt.maxx, tt.maxy, tt.col, uniform)
		}
	}
}

func TestGrayScannerIsUniform(t *testing.T) {
	ss := []string{
		"  0,   0,   0",
		"122,   0,   0",
		"  0,  24,  24",
	}

	var tests = []struct {
		minx, miny, maxx, maxy int
		col                    color.Color
		uniform                bool
	}{
		{0, 0, 3, 3, nil, false},
		{0, 1, 1, 2, color.Gray{122}, true},
		{1, 2, 3, 3, color.Gray{24}, true},
	}

	img := newGrayFromString(ss)
	scanner, err := NewScanner(img)
	test.Check(t, err)
	for _, tt := range tests {
		// FIXME: only for debug, remove it
		//sub := img.SubImage(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy))
		uniform, col := scanner.IsUniform(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy))
		if uniform != tt.uniform {
			//t.Logf("sub image %v\n", sub)
			t.Errorf("want uniform=%v for IsUniform(rect{%d,%d|%d,%d}), got %v", tt.uniform, tt.minx, tt.miny, tt.maxx, tt.maxy, uniform)
		}
		if col != tt.col {
			//t.Logf("sub image %v\n", sub)
			t.Errorf("want color=%v for IsUniform(rect{%d,%d|%d,%d}), got %v", tt.col, tt.minx, tt.miny, tt.maxx, tt.maxy, col)
		}
	}
}

func TestGrayScannerAverageColor(t *testing.T) {
	ss := []string{
		"  0,   0,   0",
		"122,   0,   0",
		"  0,  24,  24",
	}

	var tests = []struct {
		minx, miny, maxx, maxy int
		col                    color.Color
		uniform                bool
	}{

		{0, 0, 1, 1, color.Gray{0}, true},
		{0, 0, 3, 3, color.Gray{18}, false},
		{0, 1, 1, 2, color.Gray{122}, true},
		{1, 2, 3, 3, color.Gray{24}, true},
	}

	scanner, err := NewScanner(newGrayFromString(ss))
	test.Check(t, err)
	for _, tt := range tests {
		uniform, col := scanner.AverageColor(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy))
		if uniform != tt.uniform {
			t.Errorf("want uniform=%v for AverageColor(rect{%d,%d|%d,%d}), got %v", tt.uniform, tt.minx, tt.miny, tt.maxx, tt.maxy, uniform)
		}
		if col != tt.col {
			t.Errorf("want color=%v for AverageColor(rect{%d,%d|%d,%d}), got %v", tt.col, tt.minx, tt.miny, tt.maxx, tt.maxy, col)
		}
	}
}
