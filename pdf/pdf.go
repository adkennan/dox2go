/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */
package pdf

import (
	"fmt"
	"io"
)

const (
	smallBuf int = iota
	largeBuf
)

const smallBufSize int = 0x2000
const largeBufSize int = 0x8000

///////////////////////////////////////////////////////////

// pdfObj is the base interface implemented by all
// PDF objects.
type pdfObj interface {
	io.WriterTo
	Id() int
	Type() string
}

func startObj(o pdfObj, w io.Writer) (n int64, err error) {
	n2, err := fmt.Fprintf(w, "%d 0 obj\r\n", o.Id())
	return int64(n2), err
}

func endObj(o pdfObj, w io.Writer) (n int64, err error) {
	n2, err := fmt.Fprint(w, "endobj\r\n")
	return int64(n2), err
}

func startStream(w io.Writer) (n int64, err error) {
	n2, err := fmt.Fprint(w, "stream\r\n")
	return int64(n2), err
}

func endStream(w io.Writer) (n int64, err error) {
	n2, err := fmt.Fprint(w, "\r\nendstream\r\n")
	return int64(n2), err
}

type pdfStructure interface {
	Writer() io.Writer
	BytesWritten() int64
	Error() error
	HandleResult(n int64, err error)

	Name(name string)
	Value(val interface{})
	Ref(o pdfObj)
	Start()
	End()
}

func psName(ps pdfStructure, name string) {
	if ps.Error() == nil {
		n, err := fmt.Fprintf(ps.Writer(), " /%s ", name)
		ps.HandleResult(int64(n), err)
	}
}

func psValue(ps pdfStructure, val interface{}) {
	if ps.Error() == nil {
		n, err := fmt.Fprint(ps.Writer(), val)
		ps.HandleResult(int64(n), err)
	}
}

func psRef(ps pdfStructure, o pdfObj) {
	if ps.Error() == nil {
		n, err := fmt.Fprintf(ps.Writer(), " %d 0 R", o.Id())
		ps.HandleResult(int64(n), err)
	}
}

type dictionaryWriter struct {
	w   io.Writer
	n   int64
	err error
}

func (self *dictionaryWriter) Writer() io.Writer {
	return self.w
}

func (self *dictionaryWriter) BytesWritten() int64 {
	return self.n
}

func (self *dictionaryWriter) Error() error {
	return self.err
}

func (self *dictionaryWriter) HandleResult(n int64, err error) {
	if self.err == nil {
		self.n += n
		self.err = err
	}
}

func (self *dictionaryWriter) Name(name string) {
	psName(self, name)
}

func (self *dictionaryWriter) Value(val interface{}) {
	psValue(self, val)
}

func (self *dictionaryWriter) Ref(o pdfObj) {
	psRef(self, o)
}

func (self *dictionaryWriter) Start() {
	psValue(self, " << ")
}

func (self *dictionaryWriter) End() {
	psValue(self, " >>\r\n")
}

type arrayWriter struct {
	w   io.Writer
	n   int64
	err error
}

func (self *arrayWriter) Writer() io.Writer {
	return self.w
}

func (self *arrayWriter) BytesWritten() int64 {
	return self.n
}

func (self *arrayWriter) Error() error {
	return self.err
}

func (self *arrayWriter) HandleResult(n int64, err error) {
	if self.err == nil {
		self.n += n
		self.err = err
	}
}

func (self *arrayWriter) Name(name string) {
	psName(self, name)
}

func (self *arrayWriter) Value(val interface{}) {
	psValue(self, val)
}

func (self *arrayWriter) Ref(o pdfObj) {
	psRef(self, o)
}

func (self *arrayWriter) Start() {
	psValue(self, " [ ")
}

func (self *arrayWriter) End() {
	psValue(self, " ]\r\n")
}
