/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */
package pdf

import (
	"bytes"
	"dox2go"
	"fmt"
	"image"
	"io"
)

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
	n, err = startObj(c, w)
	if err != nil {
		return 0, err
	}

	dw := dictionaryWriter{w, 0, nil}
	dw.Start()
	dw.Name("Type")
	dw.Name(c.Type())

	for _, o := range c.objs {
		dw.Name(o.Type())
		dw.Ref(o)
	}

	dw.End()

	if dw.err != nil {
		return n, dw.err
	}
	n += dw.n

	n2, err := endObj(c, w)
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
	n, err = startObj(c, w)
	if err != nil {
		return 0, err
	}

	dw := dictionaryWriter{w, 0, nil}
	aw := arrayWriter{w, 0, nil}
	dw.Start()
	dw.Name("Type")
	dw.Name(c.Type())
	dw.Name("Count")
	dw.Value(len(c.pages))
	dw.Name("Kids")
	aw.Start()
	for _, p := range c.pages {
		aw.Ref(p)
	}
	aw.End()
	dw.End()

	if dw.err != nil {
		return n, dw.err
	}
	if aw.err != nil {
		return n, aw.err
	}
	n = n + dw.n + aw.n

	n2, err := endObj(c, w)
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
	n, err = startObj(c, w)
	if err != nil {
		return 0, err
	}

	dw := dictionaryWriter{w, 0, nil}
	dw.Start()
	dw.Name("Type")
	dw.Name(c.Type())
	dw.Name("Count")
	dw.Value(len(c.objs))
	for _, o := range c.objs {
		dw.Name(o.Type())
		dw.Ref(o)
	}
	dw.End()

	if dw.err != nil {
		return n, dw.err
	}
	n += dw.n

	n2, err := endObj(c, w)
	n += int64(n2)
	return n, err
}

///////////////////////////////////////////////////////////

type pdfProcSet struct {
	id    int
	names []string
}

func (p *pdfProcSet) Id() int {
	return p.id
}

func (p *pdfProcSet) Type() string {
	return "ProSet"
}

func (p *pdfProcSet) WriteTo(w io.Writer) (n int64, err error) {

	n, err = startObj(p, w)
	if err != nil {
		return 0, err
	}

	aw := arrayWriter{w, 0, nil}
	aw.Start()
	for _, name := range p.names {
		aw.Name(name)
	}
	aw.End()

	if aw.err != nil {
		return n, aw.err
	}
	n += aw.n

	n2, err := endObj(p, w)
	n += int64(n2)
	return n, err
}

func (p *pdfProcSet) Add(name string) {
	p.names = append(p.names, name)
}

///////////////////////////////////////////////////////////

type pdfDoc struct {
	w        io.Writer
	pool     *dox2go.BufferPool
	objs     []pdfObj
	catalog  *pdfCatalog
	outlines *pdfOutlines
	pages    *pdfPages
	procSet  *pdfProcSet
	fonts    pdfTypeFaceList
}

func NewPdfDoc(w io.Writer) dox2go.Document {

	cat := &pdfCatalog{1, make([]pdfObj, 0, 4)}

	outlines := &pdfOutlines{2, make([]pdfObj, 0)}
	pages := &pdfPages{3, make([]*pdfPage, 0, 4)}
	procSet := &pdfProcSet{4, make([]string, 0, 2)}

	cat.objs = append(cat.objs, outlines, pages)

	bp := dox2go.NewBufferPool()
	bp.CreateCategory(smallBuf, smallBufSize)
	bp.CreateCategory(largeBuf, largeBufSize)

	doc := &pdfDoc{
		w,
		bp,
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

func (doc *pdfDoc) CreatePage(pu dox2go.Unit, size dox2go.Point, po dox2go.PageOrientation) dox2go.Page {
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

func (doc *pdfDoc) CreateFont(name string, fs dox2go.FontStyle, size float64) dox2go.Font {

	tf := doc.fonts.findTypeFace(name, fs)
	if tf == nil {
		tf = newTypeFace(len(doc.objs)+1, name, fs)
		doc.objs = append(doc.objs, tf)
		doc.fonts = append(doc.fonts, tf)
	}

	return &pdfFont{tf, size}
}

func (doc *pdfDoc) CreateImage(src image.Image) dox2go.Image {

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

	xref := doc.pool.GetBuffer(smallBuf)
	defer doc.pool.FreeBuffer(xref)

	n, err := fmt.Fprintf(xref, "xref\r\n0 %d\r\n", len(doc.objs)+1)
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

	_, err = fmt.Fprint(xref, "0000000000 65535 f\r\n")
	if err != nil {
		return err
	}

	for _, r := range xrefs {
		err = writeXrefEntry(xref, r)
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

	dw := dictionaryWriter{doc.w, 0, nil}
	dw.Start()
	dw.Name("Size")
	dw.Value(len(doc.objs) + 1)
	dw.Name("Root")
	dw.Ref(doc.catalog)
	dw.End()

	if dw.err != nil {
		return dw.err
	}

	_, err = fmt.Fprintf(doc.w, "startxref\r\n%d\r\n%%%%EOF\r\n", offset)

	return err
}
