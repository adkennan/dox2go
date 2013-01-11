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

type pdfPage struct {
	id     int
	size   dox2go.Point
	po     dox2go.PageOrientation
	pu     dox2go.Unit
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
	n, err = startObj(p, w)
	if err != nil {
		return 0, err
	}

	var width float64
	var height float64

	if p.po == dox2go.PO_Portrait {
		width = p.size.X
		height = p.size.Y
	} else {
		width = p.size.Y
		height = p.size.X
	}

	dw := dictionaryWriter{w, 0, nil}
	aw := arrayWriter{w, 0, nil}

	dw.Start()
	dw.Name("Type")
	dw.Name(p.Type())
	dw.Name("Parent")
	dw.Ref(p.parent)
	dw.Name("MediaBox")
	aw.Start()
	aw.Value(0)
	aw.Value(0)
	aw.Value(dox2go.ConvertUnit(width, p.pu, dox2go.U_PT))
	aw.Value(dox2go.ConvertUnit(height, p.pu, dox2go.U_PT))
	aw.End()
	dw.Name("Contents")
	dw.Ref(p.c)
	dw.Name("Resources")
	dw.Start()
	dw.Name("Font")
	dw.Start()
	for _, f := range p.sfc.fonts {
		dw.Value("/F")
		dw.Value(f.Id())
		dw.Value(" ")
		dw.Ref(f)
	}
	dw.End()
	dw.Name("XObject")
	dw.Start()
	for key, xo := range p.sfc.xobjs {
		dw.Name(key)
		dw.Ref(xo)
	}
	dw.End()
	dw.End()
	dw.End()

	if dw.err != nil {
		return n, dw.err
	}
	if aw.err != nil {
		return n, aw.err
	}
	n = n + dw.n + aw.n

	n2, err := endObj(p, w)
	n += int64(n2)
	return n, err
}

func (p *pdfPage) Surface() dox2go.Surface {
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
