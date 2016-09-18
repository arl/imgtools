package binimg

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"testing"
)

type Converted struct {
	Img image.Image
	Mod color.Model
}

// We return the new color model...
func (c *Converted) ColorModel() color.Model {
	return c.Mod
}

// ... but the original bounds
func (c *Converted) Bounds() image.Rectangle {
	return c.Img.Bounds()
}

// At forwards the call to the original image and
// then asks the color model to convert it.
func (c *Converted) At(x, y int) color.Color {
	return c.Mod.Convert(c.Img.At(x, y))
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestColorModel(t *testing.T) {
	err := loadPNG2("./testdata/hr_giger_13.png")
	check(t, err)
}

func TestConvertFromImage(t *testing.T) {
	src, err := loadPNG("./testdata/hr_giger_13.png")
	check(t, err)

	dst := NewFromImage(src)
	savePNG(dst, "converted.png")
}

func savePNG(img image.Image, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		return err
	}
	return nil
}

func loadPNG(filename string) (image.Image, error) {
	var (
		f   *os.File
		img image.Image
		err error
	)

	f, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err = png.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

//func loadPNG(filename string) (*Binary, error) {
func loadPNG2(filename string) error {
	var (
		f   *os.File
		img image.Image
		err error
	)

	f, err = os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err = png.Decode(f)
	if err != nil {
		return err
	}

	// Since Converted implements image, this is now a grayscale image
	bin := &Converted{img, BinaryModelHighThreshold}

	if bin.ColorModel() != BinaryModelHighThreshold {
		panic("Expecting to be able to compare color model")
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

	outfile, err := os.Create("out.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()

	png.Encode(outfile, bin)

	return nil
}
