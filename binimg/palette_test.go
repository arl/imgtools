package binimg

import (
	"image/color"
	"testing"
)

func TestPaletteBitConvert(t *testing.T) {
	var testTbl = []struct {
		pal  Palette     // palette to use
		col  color.Color // color to convert
		want Bit         // converted Bit
	}{
		{BlackAndWhite, color.RGBA{235, 152, 30, 0}, On},
		{BlackAndWhite, color.RGBA{152, 206, 39, 0}, On},
		{BlackAndWhite, color.RGBA{152, 206, 39, 255}, On},
		{BlackAndWhite, color.RGBA{228, 75, 38, 0}, On},
		{BlackAndWhite, color.RGBA{9, 145, 210, 0}, On},
		{BlackAndWhite, color.RGBA{56, 35, 133, 0}, Off},
		{BlackAndWhite, color.RGBA{168, 7, 122, 0}, Off},
		{BlackAndWhite, color.RGBA{168, 7, 122, 255}, Off},
		{BlackAndWhite, color.RGBA{224, 16, 56, 0}, Off},
	}

	for _, tt := range testTbl {

		bit := tt.pal.ConvertBit(tt.col)
		if bit != tt.want {
			t.Errorf("Palette %#v convert %v, got %v, want %v", tt.pal, tt.col, bit, tt.want)
		}
	}
}

func TestPaletteConvert(t *testing.T) {
	custom := Palette{97, color.RGBA{255, 23, 45, 12}, color.RGBA{235, 152, 30, 0}}
	var testTbl = []struct {
		pal  Palette     // palette to use
		col  color.Color // color to convert
		want color.Color // converted color.Color
	}{
		// predefined palette
		{BlackAndWhite, color.RGBA{235, 152, 30, 0}, color.White},
		{BlackAndWhite, color.RGBA{152, 206, 39, 0}, color.White},
		{BlackAndWhite, color.RGBA{152, 206, 39, 255}, color.White},
		{BlackAndWhite, color.RGBA{228, 75, 38, 0}, color.White},
		{BlackAndWhite, color.RGBA{9, 145, 210, 0}, color.White},
		{BlackAndWhite, color.RGBA{56, 35, 133, 0}, color.Black},
		{BlackAndWhite, color.RGBA{168, 7, 122, 0}, color.Black},
		{BlackAndWhite, color.RGBA{168, 7, 122, 255}, color.Black},
		{BlackAndWhite, color.RGBA{224, 16, 56, 0}, color.Black},

		// custom palette
		{custom, color.RGBA{152, 206, 39, 0}, custom.OnColor},
		{custom, color.RGBA{152, 206, 39, 255}, custom.OnColor},
		{custom, color.RGBA{228, 75, 38, 0}, custom.OnColor},
		{custom, color.RGBA{9, 145, 210, 0}, custom.OnColor},
		{custom, color.RGBA{56, 35, 133, 0}, custom.OffColor},
		{custom, color.RGBA{168, 7, 122, 0}, custom.OffColor},
		{custom, color.RGBA{168, 7, 122, 255}, custom.OffColor},
		{custom, color.RGBA{224, 16, 56, 0}, custom.OffColor},
	}

	for _, tt := range testTbl {

		bit := tt.pal.Convert(tt.col)
		if bit != tt.want {
			t.Errorf("Palette %#v convert %v, got %v, want %v", tt.pal, tt.col, bit, tt.want)
		}
	}
}
