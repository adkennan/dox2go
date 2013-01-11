/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */
package pdf

import (
	"bytes"
	"image"
	"io"
)

type pdfImage struct {
	id   int
	src  image.Image
	mask *pdfImageMask
}

func (i *pdfImage) Id() int {
	return i.id
}

func (i *pdfImage) Width() int {
	return i.src.Bounds().Dx()
}

func (i *pdfImage) Height() int {
	return i.src.Bounds().Dy()
}

func (i *pdfImage) Type() string {
	return "XObject"
}

func (i *pdfImage) WriteTo(w io.Writer) (n int64, err error) {
	n, err = startObj(i, w)
	if err != nil {
		return 0, err
	}

	dw := dictionaryWriter{w, 0, nil}
	dw.Start()
	dw.Name("Type")
	dw.Name(i.Type())
	dw.Name("Subtype")
	dw.Name("Image")
	dw.Name("ColorSpace")
	dw.Name("DeviceRGB")
	dw.Name("BitsPerComponent")
	dw.Value(8)
	dw.Name("Width")
	dw.Value(i.Width())
	dw.Name("Height")
	dw.Value(i.Height())
	dw.Name("Length")
	dw.Value(3 * i.Width() * i.Height())
	dw.Name("SMask")
	dw.Ref(i.mask)
	dw.End()

	if dw.err != nil {
		return 0, dw.err
	}
	n += dw.n

	n2, err := startStream(w)
	if err != nil {
		return n, err
	}
	n += int64(n2)

	b := i.src.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := i.src.At(x, y).RGBA()
			n3, err := w.Write([]byte{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)})
			if err != nil {
				return n, err
			}
			i.mask.content.WriteByte(uint8(a >> 8))
			n += int64(n3)
		}
	}

	i.mask.w = i.Width()
	i.mask.h = i.Height()

	n2, err = endStream(w)
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = endObj(i, w)
	n += int64(n2)
	return n, err
}

type pdfImageMask struct {
	id      int
	w       int
	h       int
	content bytes.Buffer
}

func (i *pdfImageMask) Id() int {
	return i.id
}

func (i *pdfImageMask) Type() string {
	return "XObject"
}

func (i *pdfImageMask) WriteTo(w io.Writer) (n int64, err error) {

	n, err = startObj(i, w)
	if err != nil {
		return 0, err
	}

	dw := dictionaryWriter{w, 0, nil}
	dw.Start()
	dw.Name("Type")
	dw.Name(i.Type())
	dw.Name("Subtype")
	dw.Name("Image")
	dw.Name("ColorSpace")
	dw.Name("DeviceGray")
	dw.Name("BitsPerComponent")
	dw.Value(8)
	dw.Name("Width")
	dw.Value(i.w)
	dw.Name("Height")
	dw.Value(i.h)
	dw.Name("Length")
	dw.Value(i.w * i.h)
	dw.End()

	if dw.err != nil {
		return 0, dw.err
	}
	n += dw.n

	n2, err := startStream(w)
	if err != nil {
		return n, err
	}
	n += n2

	n3, err := i.content.WriteTo(w)
	if err != nil {
		return n, err
	}
	n += int64(n3)

	n2, err = endStream(w)
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = endObj(i, w)
	n += int64(n2)
	return n, err
}
