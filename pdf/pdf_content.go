/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */
package pdf

import (
	"bytes"
	"io"
)

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
	n, err = startObj(c, w)
	if err != nil {
		return 0, err
	}

	dw := dictionaryWriter{w, 0, nil}
	dw.Start()
	dw.Name("Length")
	dw.Value(c.b.Len())
	dw.End()

	if dw.err != nil {
		return n, dw.err
	}
	n += dw.n

	n2, err := startStream(w)
	if err != nil {
		return n, err
	}
	n += n2
	n2, err = c.b.WriteTo(w)
	if err != nil {
		return n, err
	}
	n += n2
	n2, err = endStream(w)
	if err != nil {
		return n, err
	}
	n += int64(n2)
	n2, err = endObj(c, w)

	n += int64(n2)
	return n, err
}
