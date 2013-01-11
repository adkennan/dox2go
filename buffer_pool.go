/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */
package dox2go

import (
	"bytes"
	"io"
)

type Buffer interface {
	Write(p []byte) (n int, err error)
	WriteTo(w io.Writer) (n int64, err error)
}

type buffer struct {
	reserved bool
	buf      *bytes.Buffer
}

func (b *buffer) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

func (b *buffer) WriteTo(w io.Writer) (n int64, err error) {
	return b.buf.WriteTo(w)
}

type pool struct {
	minSize int
	buffers []*buffer
}

func (p *pool) getBuffer() Buffer {

	for _, b := range p.buffers {
		if !b.reserved {
			b.reserved = true
			b.buf.Reset()
			return b
		}
	}

	b := &buffer{true, bytes.NewBuffer(make([]byte, 0, p.minSize))}

	p.buffers = append(p.buffers, b)

	return b
}

type BufferPool struct {
	pools map[int]*pool
}

func (bp *BufferPool) CreateCategory(category int, minBufferSize int) {
	bp.pools[category] = &pool{minBufferSize, make([]*buffer, 0, 1)}
}

func (bp *BufferPool) GetBuffer(category int) Buffer {

	return bp.pools[category].getBuffer()
}

func (bp *BufferPool) FreeBuffer(w Buffer) {

	if b, ok := w.(*buffer); ok {

		b.reserved = false

	} else {

		panic("Not a pool buffer!")
	}
}

func NewBufferPool() *BufferPool {
	return &BufferPool{make(map[int]*pool)}
}
