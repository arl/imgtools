# imgtools

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/aurelien-rainone/imgtools) [![Build Status](https://travis-ci.org/aurelien-rainone/imgtools.svg?branch=master)](https://travis-ci.org/aurelien-rainone/imgtools) [![Coverage Status](https://coveralls.io/repos/github/aurelien-rainone/imgtools/badge.svg?branch=master)](https://coveralls.io/github/aurelien-rainone/imgtools?branch=master)


`imgtools` package contains some utilities for working with 2D images in Go,
that complete the standard Go `image` package.

- `imgtools/binimg`: a binary image format, that is an image that has only two
possible values for each pixel. Typically, the two colors used for a binary
image are black and white, though any two colors can be used. Such images are
also referred to as *bi-level*, or *two-level*. Each pixel is stored as a
single bit, 0 represents the OffColor and 1 the OnColor in the two colors
Palette contained in the image.

- `imgscan` sub-package
