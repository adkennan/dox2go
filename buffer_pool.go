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

// Buffer us a reusable block of memory used as a temporary
// workspace.
// 
// Write writes a slice of bytes to the buffer.
//
// WriteTo writes the contents of the buffer to the supplied
// io.Writer.
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

// BufferPool maintains a collection of reusable Buffers.
// 
// Callers request a buffer of a given category and then
// return it to the pool once they are done with it.
type BufferPool struct {
	pools map[int]*pool
}

// CreateCategory adds a category of buffer such as "small" or "large".
func (bp *BufferPool) CreateCategory(category int, minBufferSize int) {
	bp.pools[category] = &pool{minBufferSize, make([]*buffer, 0, 1)}
}

// GetBuffer returns a buffer of the requested category. If no free
// buffers of the requested categories are available a new buffer 
// will be allocated.
func (bp *BufferPool) GetBuffer(category int) Buffer {

	return bp.pools[category].getBuffer()
}

// FreeBuffer returns a buffer to the pool, allowing other callers to use it.
func (bp *BufferPool) FreeBuffer(w Buffer) {

	if b, ok := w.(*buffer); ok {

		b.reserved = false

	} else {

		panic("Not a pool buffer!")
	}
}

// NewBufferPool constructs a new buffer pool.
func NewBufferPool() *BufferPool {
	return &BufferPool{make(map[int]*pool)}
}
