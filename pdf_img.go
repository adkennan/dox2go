package dox2go

import (
	"bytes"
	"fmt"
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
	n2, err := startObj(i, w)
	if err != nil {
		return 0, err
	}
	n = int64(n2)
	d := pdfDict{"Type": pn(i.Type()),
		"Subtype":          "/Image",
		"ColorSpace":       "/DeviceRGB",
		"BitsPerComponent": 8,
		"Width":            i.Width(),
		"Height":           i.Height(),
		"Length":           3 * i.Width() * i.Height(),
		"SMask":             refObj(i.mask)}

	n2, err = fmt.Fprint(w, d)
	if err != nil {
		return n, err
	}
	n += int64(n2)

	n2, err = fmt.Fprint(w, "stream\r\n")
	if err != nil {
		return n, err
	}
	n += int64(n2)

	b := i.src.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := i.src.At(x, y).RGBA()
			n2, err = w.Write([]byte{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)})
			if err != nil {
				return n, err
			}
			i.mask.content.WriteByte(uint8(a >> 8))
			n += int64(n2)
		}
	}

	i.mask.w = i.Width()
	i.mask.h = i.Height()

	n2, err = fmt.Fprint(w, "\r\nendstream\r\n")
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

	n2, err := startObj(i, w)
	if err != nil {
		return 0, err
	}
	n = int64(n2)
	d := pdfDict{"Type": pn(i.Type()),
		"Subtype":          "/Image",
		"ColorSpace":       "/DeviceGray",
		"BitsPerComponent": 8,
		"Width":            i.w,
		"Height":           i.h,
		"Length":           i.w * i.h}

	n2, err = fmt.Fprint(w, d)
	if err != nil {
		return n, err
	}
	n += int64(n2)

	n2, err = fmt.Fprint(w, "stream\r\n")
	if err != nil {
		return n, err
	}
	n += int64(n2)

	n3, err := i.content.WriteTo(w)
	if err != nil {
		return n, err
	}
	n += int64(n3)

	n2, err = fmt.Fprint(w, "\r\nendstream\r\n")
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = endObj(i, w)
	n += int64(n2)
	return n, err
}
