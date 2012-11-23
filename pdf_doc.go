/*
* dox2go - A document generating library for go.
*
* Copyright 2011 Andrew Kennan. All rights reserved.
*
* Redistribution and use in source and binary forms, with or without modification, are
* permitted provided that the following conditions are met:
*
* 1. Redistributions of source code must retain the above copyright notice, this list of
* conditions and the following disclaimer.
*
* 2. Redistributions in binary form must reproduce the above copyright notice, this list
* of conditions and the following disclaimer in the documentation and/or other materials
* provided with the distribution.
*
* THIS SOFTWARE IS PROVIDED BY ANDREW KENNAN ''AS IS'' AND ANY EXPRESS OR IMPLIED
* WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
* FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL ANDREW KENNAN OR
* CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
* CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
* SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
* ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
* NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
* ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*
* The views and conclusions contained in the software and documentation are those of the
* authors and should not be interpreted as representing official policies, either expressed
* or implied, of Andrew Kennan.
 */
package dox2go

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"strconv"
)

///////////////////////////////////////////////////////////

// pdfObj is the base interface implemented by all
// PDF objects.
type pdfObj interface {
	io.WriterTo
	Id() int
	Type() string
}

func startObj(o pdfObj, w io.Writer) (n int, err error) {
	return fmt.Fprintf(w, "%d 0 obj\r\n", o.Id())
}

func endObj(o pdfObj, w io.Writer) (n int, err error) {
	return fmt.Fprint(w, "endobj\r\n")
}

func refObj(o pdfObj) string {
	return strconv.Itoa(o.Id()) + " 0 R"
}

///////////////////////////////////////////////////////////

type pdfName struct {
	name string
}

func (n *pdfName) String() string {
	return "/" + n.name
}

func pn(n string) *pdfName {
	return &pdfName{n}
}

///////////////////////////////////////////////////////////

type pdfDict map[string]interface{}

func (d pdfDict) String() string {
	var b bytes.Buffer

	fmt.Fprint(&b, " <<")

	for key, value := range d {
		fmt.Fprintf(&b, " /%v %v\r\n", key, value)
	}

	fmt.Fprint(&b, " >>\r\n")

	return b.String()
}

///////////////////////////////////////////////////////////

type pdfArray []interface{}

func (a pdfArray) String() string {
	var b bytes.Buffer

	fmt.Fprint(&b, "[")
	for _, el := range a {
		fmt.Fprintf(&b, "%v ", el)
	}
	fmt.Fprint(&b, "]")

	return b.String()
}

///////////////////////////////////////////////////////////

type pdfPage struct {
	id     int
	size   Point
	po     PageOrientation
	pu     Unit
	parent pdfObj
	sfc    *pdfSurface
	c      *pdfContent
}

func (p *pdfPage) Id() int {
	return p.id
}

func (p *pdfPage) Close() {
	if p.sfc != nil {
		p.sfc.Close()
	}
}

func (p *pdfPage) Type() string {
	return "Page"
}

func (p *pdfPage) WriteTo(w io.Writer) (n int64, err error) {
	n2, err := startObj(p, w)
	if err != nil {
		return 0, err
	}
	n = int64(n2)

	fs := make(pdfDict)
	for _, f := range p.sfc.fonts {
		fs["F"+strconv.Itoa(f.Id())] = refObj(f)
	}

	xos := make(pdfDict)
	for key, xo := range p.sfc.xobjs {
		xos[key] = refObj(xo)
	}

	var width float64
	var height float64

	if p.po == PO_Portrait {
		width = p.size.X
		height = p.size.Y
	} else {
		width = p.size.Y
		height = p.size.X
	}

	n2, err = fmt.Fprint(w, pdfDict{
		"Type":   pn(p.Type()),
		"Parent": refObj(p.parent),
		"MediaBox": &pdfArray{0, 0,
			ConvertUnit(width, p.pu, U_PT),
			ConvertUnit(height, p.pu, U_PT)},
		"Contents": refObj(p.c),
		"Resources": pdfDict{
			"Font":    fs,
			"XObject": xos,
		},
	})
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = endObj(p, w)
	n += int64(n2)
	return n, err
}

func (p *pdfPage) Surface() Surface {
	if p.sfc == nil {
		p.sfc = &pdfSurface{
			p.c.b,
			p.pu,
			false,
			make([]*pdfTypeFace, 0, 4),
			nil,
			make(map[string]pdfObj),
		}
	}
	return p.sfc
}

///////////////////////////////////////////////////////////

type pdfContent struct {
	id int
	b  *bytes.Buffer
}

func (c *pdfContent) Id() int {
	return c.id
}

func (p *pdfContent) Type() string {
	return "Content"
}

func (c *pdfContent) WriteTo(w io.Writer) (n int64, err error) {
	n2, err := startObj(c, w)
	if err != nil {
		return 0, err
	}
	n = int64(n2)
	n2, err = fmt.Fprint(w, pdfDict{
		"Length": c.b.Len(),
	})
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = fmt.Fprint(w, "stream\r\n")
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n3, err := c.b.WriteTo(w)
	if err != nil {
		return n, err
	}
	n += n3
	n2, err = fmt.Fprint(w, "\r\nendstream\r\n")
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = endObj(c, w)
	n += int64(n2)
	return n, err
}

///////////////////////////////////////////////////////////

type pdfCatalog struct {
	id   int
	objs []pdfObj
}

func (c *pdfCatalog) Id() int {
	return c.id
}

func (p *pdfCatalog) Type() string {
	return "Catalog"
}

func (c *pdfCatalog) WriteTo(w io.Writer) (n int64, err error) {
	n2, err := startObj(c, w)
	if err != nil {
		return 0, err
	}
	n = int64(n2)
	d := pdfDict{"Type": pn(c.Type())}

	for _, o := range c.objs {
		d[o.Type()] = refObj(o)
	}

	n2, err = fmt.Fprint(w, d)
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = endObj(c, w)
	n += int64(n2)
	return n, err
}

///////////////////////////////////////////////////////////

type pdfPages struct {
	id    int
	pages []*pdfPage
}

func (c *pdfPages) Id() int {
	return c.id
}

func (p *pdfPages) Type() string {
	return "Pages"
}

func (c *pdfPages) WriteTo(w io.Writer) (n int64, err error) {
	n2, err := startObj(c, w)
	if err != nil {
		return 0, err
	}
	n = int64(n2)
	l := len(c.pages)
	d := pdfDict{"Type": pn(c.Type()), "Count": l}

	k := make(pdfArray, l, l)
	for ix, p := range c.pages {
		k[ix] = refObj(p) + "\r\n"
	}
	d["Kids"] = k.String()

	n2, err = fmt.Fprint(w, d)
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = endObj(c, w)
	n += int64(n2)
	return n, err
}

func (c *pdfPages) Close() {
	for _, p := range c.pages {
		p.Close()
	}
}

///////////////////////////////////////////////////////////

type pdfOutlines struct {
	id   int
	objs []pdfObj
}

func (c *pdfOutlines) Id() int {
	return c.id
}

func (p *pdfOutlines) Type() string {
	return "Outlines"
}

func (c *pdfOutlines) WriteTo(w io.Writer) (n int64, err error) {
	n2, err := startObj(c, w)
	if err != nil {
		return 0, err
	}
	n = int64(n2)
	d := pdfDict{"Type": pn(c.Type()), "Count": len(c.objs)}

	for _, o := range c.objs {
		d[o.Type()] = refObj(o)
	}

	n2, err = fmt.Fprint(w, d)
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = endObj(c, w)
	n += int64(n2)
	return n, err
}

///////////////////////////////////////////////////////////

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

func newTypeFace(id int, name string, fs FontStyle) *pdfTypeFace {
	return &pdfTypeFace{
		id,
		fst_Type1,
		defaultFaces[defaultFaceIndex(name)][fs],
		fs,
	}
}

type pdfTypeFaceList []*pdfTypeFace

func (fonts pdfTypeFaceList) findTypeFace(name string, fs FontStyle) *pdfTypeFace {
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
	fs       FontStyle
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

	n2, err := startObj(f, w)
	if err != nil {
		return 0, err
	}
	n = int64(n2)

	n2, err = fmt.Fprint(w, pdfDict{
		"Type":     pn(f.Type()),
		"Subtype":  pn(fontTypeString(f.subType)),
		"BaseFont": pn(f.baseFont),
		"Name":     pn("F" + strconv.Itoa(f.id)),
	})
	if err != nil {
		return n, err
	}

	n += int64(n2)
	n2, err = endObj(f, w)
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

func (f *pdfFont) Style() FontStyle {
	return f.face.fs
}

func (f *pdfFont) Size() float64 {
	return f.size
}

func (f *pdfFont) Equals(other *pdfFont) bool {
	return f.face.id == other.face.id &&
		f.size == other.size
}

///////////////////////////////////////////////////////////

type pdfProcSet struct {
	id    int
	names pdfArray
}

func (p *pdfProcSet) Id() int {
	return p.id
}

func (p *pdfProcSet) Type() string {
	return "ProSet"
}

func (p *pdfProcSet) WriteTo(w io.Writer) (n int64, err error) {

	n2, err := startObj(p, w)
	if err != nil {
		return 0, err
	}
	n = int64(n2)

	n2, err = fmt.Fprintf(w, "%s\r\n", p.names)
	if err != nil {
		return 0, err
	}
	n = int64(n2)

	n2, err = endObj(p, w)
	n += int64(n2)
	return n, err
}

func (p *pdfProcSet) Add(name string) {
	p.names = append(p.names, pn(name))
}

///////////////////////////////////////////////////////////

type pdfDoc struct {
	w        io.Writer
	objs     []pdfObj
	catalog  *pdfCatalog
	outlines *pdfOutlines
	pages    *pdfPages
	procSet  *pdfProcSet
	fonts    pdfTypeFaceList
}

func NewPdfDoc(w io.Writer) Document {

	cat := &pdfCatalog{1, make([]pdfObj, 0, 4)}

	outlines := &pdfOutlines{2, make([]pdfObj, 0)}
	pages := &pdfPages{3, make([]*pdfPage, 0, 4)}
	procSet := &pdfProcSet{4, make(pdfArray, 0, 2)}

	cat.objs = append(cat.objs, outlines, pages)

	doc := &pdfDoc{
		w,
		make([]pdfObj, 0, 10),
		cat,
		outlines,
		pages,
		procSet,
		make([]*pdfTypeFace, 0, 4),
	}

	doc.objs = append(doc.objs, cat, outlines, pages, procSet)

	procSet.Add("PDF")
	procSet.Add("Text")

	return doc
}

func (doc *pdfDoc) CreatePage(pu Unit, size Point, po PageOrientation) Page {
	p := &pdfPage{
		len(doc.objs) + 1,
		size,
		po,
		pu,
		doc.pages,
		nil,
		&pdfContent{
			len(doc.objs) + 2,
			new(bytes.Buffer),
		},
	}

	doc.objs = append(doc.objs, p, p.c)
	doc.pages.pages = append(doc.pages.pages, p)

	return p
}

func (doc *pdfDoc) CreateFont(name string, fs FontStyle, size float64) Font {

	tf := doc.fonts.findTypeFace(name, fs)
	if tf == nil {
		tf = newTypeFace(len(doc.objs)+1, name, fs)
		doc.objs = append(doc.objs, tf)
		doc.fonts = append(doc.fonts, tf)
	}

	return &pdfFont{tf, size}
}

func (doc *pdfDoc) CreateImage(src image.Image) Image {

	m := &pdfImageMask{len(doc.objs) + 2, 0, 0, bytes.Buffer{}}
	i := &pdfImage{len(doc.objs) + 1, src, m}

	doc.objs = append(doc.objs, i, m)

	return i
}

func writeXrefEntry(w io.Writer, offset int64) (err error) {
	_, err = fmt.Fprintf(w, "%010d 00000 n\r\n", offset)
	return err
}

func (doc *pdfDoc) Close() (err error) {

	doc.pages.Close()

	var offset int64 = 0

	var xref bytes.Buffer
	n, err := fmt.Fprintf(&xref, "xref\r\n0 %d\r\n", len(doc.objs)+1)
	if err != nil {
		return err
	}

	n, err = fmt.Fprint(doc.w, "%PDF-1.7\r\n")
	if err != nil {
		return err
	}
	offset += int64(n)

	n, err = doc.w.Write([]byte{200, 200, 200, 200, 13, 10})
	if err != nil {
		return err
	}
	offset += int64(n)

	xrefs := make([]int64, len(doc.objs))
	for ix, o := range doc.objs {
		xrefs[ix] = offset
		n2, err := o.WriteTo(doc.w)
		if err != nil {
			return err
		}
		offset += n2
	}

	_, err = fmt.Fprint(&xref, "0000000000 65535 f\r\n")
	if err != nil {
		return err
	}

	for _, r := range xrefs {
		err = writeXrefEntry(&xref, r)
		if err != nil {
			return err
		}
	}

	_, err = xref.WriteTo(doc.w)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(doc.w, "trailer\r\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(doc.w, pdfDict{
		"Size": len(doc.objs) + 1,
		"Root": refObj(doc.catalog),
	})

	fmt.Fprintf(doc.w, "startxref\r\n%d\r\n%%%%EOF\r\n", offset)

	return err
}
