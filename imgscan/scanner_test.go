package imgscan

import (
	"image"
	"testing"

	"github.com/arl/imgtools/binimg"
)

func TestScannerSupportedImageTypes(t *testing.T) {
	var r = image.Rect(0, 0, 16, 16)
	var tests = []struct {
		img  image.Image // image type to be tested
		want error       // error returned by NewScanner
	}{
		{binimg.New(r), nil},
		{image.NewGray(r), nil},
		{image.NewRGBA(r), ErrUnsupportedType},
		{image.NewNRGBA(r), ErrUnsupportedType},
	}

	for _, tt := range tests {
		gotImg, gotErr := NewScanner(tt.img)

		switch tt.want {

		// image type should be supported
		case nil:
			if gotImg == nil {
				t.Errorf("want image of type %T supported by NewScanner, got img=%v", tt.img, tt.img)
			}
			if gotErr == ErrUnsupportedType {
				t.Errorf("want image of type %T supported by NewScanner, got err=%v", tt.img, gotErr)
			}

		// image type is not suppported
		case ErrUnsupportedType:
			if gotImg != nil {
				t.Errorf("want image of type %T NOT supported by NewScanner, got img != nil", tt.img)
			}
			if gotErr == nil || gotErr != ErrUnsupportedType {
				t.Errorf("want image of type %T NOT supported by NewScanner, got err=%v", tt.img, gotErr)
			}
		}
	}
}
