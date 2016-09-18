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

//continuer ici: en fait ce qu'on a demontre ici c'est qu'on peut facilement
//transformer une image monochrome (avec different thresholds, et pourquoi pas la possibilité de proposer soi-meme son BinaryModel, avec threshold et couleur customisée)
//mais bon, afin de pouvoir scanner l'image efficacement, le Scanner a besoin d'une imgbin.Binary et non d'une image.Image

//C'est pourquoi la conversion devrait se faire comme
//https: //www.socketloop.com/tutorials/golang-convert-png-transparent-background-image-to-jpg-or-jpeg-image,
//c'est à dire en creée uneNewBinary et en copiant l'image Src dedans, a
//partir de la on a une imgbin.Binary que l'on peut placer dans un scanner,
//scanner qui est optimisé pour un certain ColorModel. On vient de demontrere
//que l'on peut verifier le ColorModel d'une image, donc en fonction de cela
//on eptu scanner comme le fait l'actuel Scanner. Ainsi on aura une interface
//bien plus saine, avec certaienement quelque chose qui vaut la peine d'avoir
//son propre package

//Ainsi le package quadtree pourrait fonctionner comme il le souhaite en
//utilisant un Scanner content une image, si le colorModel de l'image est
//BinaryModel on applique bytes.IndexBytes pour le scanning, sinon un simple
//BruteForceScanner
