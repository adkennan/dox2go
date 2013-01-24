/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */

package pdf

import (
	"dox2go"
	"io"
)

// Standard PDF fonts.
const (
	FONT_Times        string = "Times"
	FONT_Helvetica           = "Helvetica"
	FONT_Courier             = "Courier"
	FONT_Symbol              = "Symbol"
	FONT_ZapfDingbats        = "ZapfDingbats"
)

var defaultFaces = [5][4]string{
	{"Times-Roman", "Times-Bold", "Times-Italic", "Times-BoldItalic"},
	{"Helvetica", "Helvetica-Bold", "Helvetica-Oblique", "Helvetica-BoldOblique"},
	{"Courier", "Courier-Bold", "Courier-Oblique", "Courier-BoldOblique"},
	{"Symbol", "Symbol", "Symbol", "Symbol"},
	{"ZapfDingbats", "ZapfDingbats", "ZapfDingbats", "ZapfDingbats"},
}

func defaultFaceIndex(name string) int {
	switch name {
	case FONT_Helvetica:
		return 1
	case FONT_Courier:
		return 2
	case FONT_Symbol:
		return 3
	case FONT_ZapfDingbats:
		return 4
	default:
		break
	}
	return 0
}

func newTypeFace(id int, name string, fs dox2go.FontStyle) *pdfTypeFace {
	return &pdfTypeFace{
		id,
		fst_Type1,
		defaultFaces[defaultFaceIndex(name)][fs],
		fs,
	}
}

type pdfTypeFaceList []*pdfTypeFace

func (fonts pdfTypeFaceList) findTypeFace(name string, fs dox2go.FontStyle) *pdfTypeFace {
	baseFont := defaultFaces[defaultFaceIndex(name)][fs]

	for _, f := range fonts {
		if f.baseFont == baseFont {
			return f
		}
	}

	return nil
}

type fontSubType int32

const (
	fst_Type1 fontSubType = iota
)

type pdfTypeFace struct {
	id       int
	subType  fontSubType
	baseFont string
	fs       dox2go.FontStyle
}

func (f *pdfTypeFace) Id() int {
	return f.id
}

func (f *pdfTypeFace) Type() string {
	return "Font"
}

func fontTypeString(ft fontSubType) string {
	switch ft {
	case fst_Type1:
		return "Type1"
	}
	return ""
}

func (f *pdfTypeFace) WriteTo(w io.Writer) (n int64, err error) {

	n, err = startObj(f, w)
	if err != nil {
		return 0, err
	}

	dw := dictionaryWriter{w, 0, nil}
	dw.Start()
	dw.Name("Type")
	dw.Name(f.Type())
	dw.Name("Subtype")
	dw.Name(fontTypeString(f.subType))
	dw.Name("BaseFont")
	dw.Name(f.baseFont)
	dw.Name("Name")
	dw.Value("F")
	dw.Value(f.id)
	dw.Value(" ")
	dw.End()

	if dw.err != nil {
		return n, dw.err
	}
	n += dw.n

	n2, err := endObj(f, w)
	n += int64(n2)
	return n, err
}

type pdfFont struct {
	face *pdfTypeFace
	size float64
}

func (f *pdfFont) Id() int {
	return f.face.id
}

func (f *pdfFont) Style() dox2go.FontStyle {
	return f.face.fs
}

func (f *pdfFont) Size() float64 {
	return f.size
}

func (f *pdfFont) Equals(other *pdfFont) bool {
	return f.face.id == other.face.id &&
		f.size == other.size
}
