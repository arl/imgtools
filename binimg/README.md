# binimg

`binimg` package proposes an in-memory binary image format, implementing the
`image.Image` interface. A binary image has only two possible colors for each
pixel, generally Black and White, though any two colors can be used.

A pixel could be stored as a single bit, but as the main goal for this package
is time-efficient manipulation of binary images, `binimg.Binary` uses `byte` as
the underlying data type for pixels.

`binimg.Binary` images are created by calling either: 

- `binimg.New()`, that creates an image of the given dimensions, filled with
`binimg.OffColor`, and providing a `binimg.Palette` or using a predefined one.
- `binimg.NewFromImage`, that takes any `image.Image` and converts it to a
`Binary`, by using the provided `binimg.Palette`.


## TODO: continue here...

-----------------------

**converted using the default color model: `binaryModel`**

<img src="../readme-images/colorgopher.png" width="128">  <img src="../readme-images/bwgopher.png" width="128">

**converted using the high threshold color model: `BinaryModelHighThreshold`**

<img src="../readme-images/colorgopher.png" width="128">  <img src="../readme-images/bwgopher.high.threshold.png" width="128">

-----------------------

## Usage

- **Create and modify new binary image**

```go
package main

import (
	"image"
	"github.com/aurelien-rainone/binimg"
)

func main() {
	// create a new image (black)
	bin := binimg.New(image.Rect(0,0, 128, 128))

	// set a pixel to White
	bin.SetBit(10, 0, White)

	// set a pixel, converting original color using the binaryModel
	bin.Set(10, 0, color.RGBA(127, 23, 798, 255))

	// set rectangular region to White
	bin.SetRect(image.Rect(32, 32, 64, 64), White)
}
```

- **Convert `image.Image` into `binimg.Binary`**

```go
package main

import "github.com/aurelien-rainone/binimg"

func main() {
	// load image ("color-gopher.png")
	// ...
	bin := binimg.NewFromImage(img)

	// save image ("black&white-gopher.png")
	// ...
}
```

- **Use a custom binary `color.Model`**

```go
package main

import "github.com/aurelien-rainone/binimg"

func main() {
	// use one of the predefined color models
	model := binimg.BinaryModelHighThreshold
	bin := binimg.NewCustomBinary(image.Rect(0, 0, 32, 32), model)

	// convert an image to black and white with a custom color model
	mymodel := binimg.NewBinaryModel(214)
	bin := binimg.NewCustomFromImage(myimg, mymodel)
}
```
