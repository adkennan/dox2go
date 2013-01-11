/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */
package dox2go

import (
	"fmt"
	"math"
)

type Point struct {
	X float64
	Y float64
}

func (p *Point) ChangeUnit(from Unit, to Unit) Point {
	return Point{
		ConvertUnit(p.X, from, to),
		ConvertUnit(p.Y, from, to)}
}

type Size struct {
	W float64
	H float64
}

func (s *Size) ChangeUnit(from Unit, to Unit) Size {
	return Size{
		ConvertUnit(s.W, from, to),
		ConvertUnit(s.H, from, to)}
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func RGB(r uint8, g uint8, b uint8) Color {
	return Color{r, g, b, 255}
}

func RGBA(r uint8, g uint8, b uint8, a uint8) Color {
	return Color{r, g, b, a}
}

type LineCapStyle int32

const (
	LC_ButtCap   LineCapStyle = 0
	LC_RoundCap               = 1
	LC_SquareCap              = 2
)

type LineJoinStyle int32

const (
	LJ_MitreJoin LineJoinStyle = iota
	LJ_RoundJoin
	LJ_BevelJoin
)

type FontStyle int32

const (
	FS_Regular FontStyle = 0
	FS_Bold    FontStyle = 1
	FS_Italic  FontStyle = 2
)

type Font interface {
	Id() int
	Style() FontStyle
	Size() float64
}

type Image interface {
	Id() int
	Width() int
	Height() int
}

type PathCmdType uint8

const (
	MoveCmdType  PathCmdType = 0xFF
	LineCmdType              = 0xFE
	CurveCmdType             = 0xFD
	RectCmdType              = 0xFC
	ArcCmdType               = 0xFB
	CloseCmdType             = 0xFA
)

type Path struct {
	elements []byte
}

func NewPath() *Path {
	return &Path{make([]byte, 0, 16)}
}

func (path *Path) ensureCap(n int) {
	if len(path.elements)+n >= cap(path.elements) {
		newEls := make([]byte, len(path.elements), len(path.elements)*2)
		copy(newEls, path.elements)
		path.elements = newEls
	}
}

func (path *Path) writeUint8(val uint8) {
	path.ensureCap(1)
	path.elements = append(path.elements, val)
}

func (path *Path) writeCmdType(val PathCmdType) {
	path.writeUint8(uint8(val))
}

func (path *Path) writeFloat64(val float64) {

	path.ensureCap(8)

	bits := math.Float64bits(val)
	for ix := 0; ix < 8; ix++ {
		path.elements = append(path.elements, byte(bits&0xFF))
		bits >>= 8
	}
}

func (path *Path) Move(to Point) {
	path.writeCmdType(MoveCmdType)
	path.writeFloat64(to.X)
	path.writeFloat64(to.Y)
}

func (path *Path) Line(to Point) {
	path.writeCmdType(LineCmdType)
	path.writeFloat64(to.X)
	path.writeFloat64(to.Y)
}

func (path *Path) Curve(control1 Point, control2 Point, to Point) {
	path.writeCmdType(CurveCmdType)
	path.writeFloat64(control1.X)
	path.writeFloat64(control1.Y)
	path.writeFloat64(control2.X)
	path.writeFloat64(control2.Y)
	path.writeFloat64(to.X)
	path.writeFloat64(to.Y)
}

func (path *Path) Rect(from Point, to Point) {
	path.writeCmdType(RectCmdType)
	path.writeFloat64(from.X)
	path.writeFloat64(from.Y)
	path.writeFloat64(to.X)
	path.writeFloat64(to.Y)
}

func (path *Path) Arc(center Point, radius, start, sweep float64) {
	path.writeCmdType(ArcCmdType)
	path.writeFloat64(center.X)
	path.writeFloat64(center.Y)
	path.writeFloat64(radius)
	path.writeFloat64(start)
	path.writeFloat64(sweep)
}

func (path *Path) Close() {
	path.writeCmdType(CloseCmdType)
}

func (path *Path) Reader() PathReader {
	return &pathReader{path.elements, 0}
}

type PathReader interface {
	ReadCommandType() (cmdType PathCmdType, ok bool)
	ReadFloat64() (val float64)
	ReadUint8() (val uint8)
	Dump()
}

type pathReader struct {
	elements []byte
	pos      int
}

func (p *pathReader) ensureCap(n int) bool {
	return p.pos+n <= len(p.elements)
}

func (p *pathReader) ReadCommandType() (cmdType PathCmdType, ok bool) {
	ok = p.ensureCap(1)
	if ok {
		cmdType = PathCmdType(p.elements[p.pos])
		p.pos++
	}
	return
}

func (p *pathReader) ReadFloat64() (val float64) {
	if !p.ensureCap(8) {
		panic("End of buffer!")
	}

	var bits uint64
	for ix := 0; ix < 8; ix++ {
		bits |= uint64(p.elements[p.pos+ix]) << uint32(ix*8)
	}

	val = math.Float64frombits(bits)
	p.pos += 8
	return
}

func (p *pathReader) ReadUint8() (val uint8) {
	if !p.ensureCap(1) {
		panic("End of buffer!")
	}

	val = p.elements[p.pos]
	p.pos++
	return val
}

func (p *pathReader) Dump() {
	for _, v := range p.elements {
		fmt.Printf("%0x ", v)
	}
	fmt.Println("")
}

type Surface interface {
	PushState()
	PopState()

	Rotate(byRadians float64)
	Skew(xRadians float64, yRadians float64)
	Translate(byDistance Point)
	Scale(xScale float64, yScale float64)

	Fg(color Color)
	Bg(color Color)
	LineWidth(width float64)
	LineCap(capStyle LineCapStyle)
	LineJoin(joinStyle LineJoinStyle)
	LinePattern(pattern []float64, phase float64)

	Text(f Font, at Point, text string)

	Image(i Image, at Point, size Size)

	Stroke(path *Path)
	Fill(path *Path)
}
