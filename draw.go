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

// Color describes a color with red, green and blue
// components as well as a transparency level.
//
// R is the red component of the color.
//
// G is the green component of the color.
//
// B is the blue component of the color.
//
// A is the alpha, or transparency, of the color.
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// RGB returns a Color based on the supplied r, g and b components.
func RGB(r uint8, g uint8, b uint8) Color {
	return Color{r, g, b, 255}
}

// RGB returns a Color based on the supplied r, g, b and transpareny
// components.
func RGBA(r uint8, g uint8, b uint8, a uint8) Color {
	return Color{r, g, b, a}
}

// LineCapStyle describes the types of line endings available.
type LineCapStyle int32

// These are the available line cap styles.
const (
	LC_ButtCap   LineCapStyle = 0
	LC_RoundCap               = 1
	LC_SquareCap              = 2
)

// LineJoinStyle describes the types of line joins available.
type LineJoinStyle int32

// These are the available line join styles.
const (
	LJ_MitreJoin LineJoinStyle = iota
	LJ_RoundJoin
	LJ_BevelJoin
)

// FontStyle defines the types of styles available for fonts.
type FontStyle int32

// The font styles available.
const (
	FS_Regular FontStyle = 0
	FS_Bold    FontStyle = 1
	FS_Italic  FontStyle = 2
)

// Font is an interface describing a typeface used for
// drawing text.
//
// Id returns the identifier assigned by the document
// to the font.
//
// Style returns the style of the font.
//
// Size returns the size of the font in the current page
// unit.
type Font interface {
	Id() int
	Style() FontStyle
	Size() float64
}

// Image is an interface that describes a bitmap that
// can be drawn on a page.
//
// Id returns am identifier assigned by the document
// to this image.
//
// Width returns the width of the image in pixels.
//
// Height returns the height of the image in pixels.
type Image interface {
	Id() int
	Width() int
	Height() int
}

// PathCmdType describes the types of drawing operations
// supported by Paths.
type PathCmdType uint8

// The types of drawing operations.
const (
	MoveCmdType  PathCmdType = 0xFF
	LineCmdType              = 0xFE
	CurveCmdType             = 0xFD
	RectCmdType              = 0xFC
	ArcCmdType               = 0xFB
	CloseCmdType             = 0xFA
)

// Path describes a series of drawing commands such
// as lines, arcs and rectangles that can be stroked
// or filled.
type Path struct {
	elements []byte
}

// NewPath constructs a new Path object.
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

func (path *Path) writeCmdType(val PathCmdType) {
	path.ensureCap(1)
	path.elements = append(path.elements, uint8(val))
}

func (path *Path) writeFloat64(val float64) {

	path.ensureCap(8)

	bits := math.Float64bits(val)
	for ix := 0; ix < 8; ix++ {
		path.elements = append(path.elements, byte(bits&0xFF))
		bits >>= 8
	}
}

// Move moves the drawing cursor to the supplied position.
func (path *Path) Move(x, y float64) {
	path.writeCmdType(MoveCmdType)
	path.writeFloat64(x)
	path.writeFloat64(y)
}

// Line draws a line from the current cursor position to 
// the new position.
func (path *Path) Line(x, y float64) {
	path.writeCmdType(LineCmdType)
	path.writeFloat64(x)
	path.writeFloat64(y)
}

// Curve adds a bezier curve to the Path. The current position
// of the drawing cursor is the start of the curve.
func (path *Path) Curve(cx1, cy1, cx2, cy2, x, y float64) {
	path.writeCmdType(CurveCmdType)
	path.writeFloat64(cx1)
	path.writeFloat64(cy1)
	path.writeFloat64(cx2)
	path.writeFloat64(cy2)
	path.writeFloat64(x)
	path.writeFloat64(y)
}

// Rect adds a rectangle to the path.
func (path *Path) Rect(x1, y1, x2, y2 float64) {
	path.writeCmdType(RectCmdType)
	path.writeFloat64(x1)
	path.writeFloat64(y1)
	path.writeFloat64(x2)
	path.writeFloat64(y2)
}

// Arc draws a an ellipse or partial ellipse centered on the supplied
// point. Start and sweep are in radians and define where the arc
// begins and how far it extends.
func (path *Path) Arc(x, y float64, radius, start, sweep float64) {
	path.writeCmdType(ArcCmdType)
	path.writeFloat64(x)
	path.writeFloat64(y)
	path.writeFloat64(radius)
	path.writeFloat64(start)
	path.writeFloat64(sweep)
}

// Close closes the path, drawing a line back to the starting position.
func (path *Path) Close() {
	path.writeCmdType(CloseCmdType)
}

// Reader returns an object used by Surface implementations to
// enumerate the drawing operations included in the path.
func (path *Path) Reader() PathReader {
	return &pathReader{path.elements, 0}
}

// PathReader is an interface that allows an implementation
// of Surface to enumerate the operations of a path.
//
// ReadCommandType reads the next operation from the path.
//
// ReadFloat64 reads a float64 from the path.
//
// ReadUint8 reads a uint8 from the path.
//
// Dump writes some debug info to stdout.
type PathReader interface {
	ReadCommandType() (cmdType PathCmdType, ok bool)
	ReadFloat64() (val float64)
	Dump()
}

type pathReader struct {
	elements []byte
	pos      int
}

func (p *pathReader) ensureCap(n int) bool {
	return p.pos+n <= len(p.elements)
}

// ReadCommandType reads the next operation from the path.
func (p *pathReader) ReadCommandType() (cmdType PathCmdType, ok bool) {
	ok = p.ensureCap(1)
	if ok {
		cmdType = PathCmdType(p.elements[p.pos])
		p.pos++
	}
	return
}

// ReadFloat64 reads a float64 from the path.
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

// Dump writes some debug info to stdout.
func (p *pathReader) Dump() {
	for _, v := range p.elements {
		fmt.Printf("%0x ", v)
	}
	fmt.Println("")
}

// Surface defines operations for drawing text and
// graphics on a page.
//
// PushState stores the current drawing state for
// later recall.
//
// PopState discards the current drawing state and
// recalls the previously pushed state.
//
// Rotate transforms the drawing surface by rotating
// drawing operations by the supplied angle.
//
// Skew transforms the drawing surface by skewing 
// drawing operations by the supplied angles.
//
// Translate transforms the drawing surface by 
// translating the operations by the supplied distances.
//
// Scale transforms the drawing surface by scaling
// operations by the supplied scales.
//
// Fg sets the color used to stroke paths.
//
// Bg sets the color used to fill paths.
//
// LineWidth sets the width of the lines uses for stroking
// paths. The width is in the current page units.
//
// LineCap sets the style of the cap added to lines of
// stroked paths.
//
// LineJoin sets the style of the joins between lines of a
// stroked path.
//
// LinePattern sets the pattern to use when stroking lines
// of a path. The pattern is in page units. Phase indicates
// where in the pattern to begin drawing.
//
// Text draws a text string on the Surface.
//
// Image draws a bitmap image on the Surface.
//
// Stroke strokes a path in the Fg color using the current 
// line width, joins and caps.
//
// Fill fills the area of a path with the current Bg color.
type Surface interface {
	PushState()
	PopState()

	Rotate(byRadians float64)
	Skew(xRadians float64, yRadians float64)
	Translate(x, y float64)
	Scale(xScale float64, yScale float64)

	Fg(color Color)
	Bg(color Color)
	LineWidth(width float64)
	LineCap(capStyle LineCapStyle)
	LineJoin(joinStyle LineJoinStyle)
	LinePattern(pattern []float64, phase float64)

	Text(f Font, x, y float64, text string)

	Image(i Image, x, y, w, h float64)

	Stroke(path *Path)

	Fill(path *Path)
}
