/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */

package pdf

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"

	d2g "github.com/adkennan/dox2go"
)

type pdfSurface struct {
	w        io.Writer
	u        d2g.PageUnit
	inText   bool
	fonts    []*pdfTypeFace
	lastFont *pdfFont
	xobjs    map[string]pdfObj
}

func (sfc *pdfSurface) addXObj(o pdfObj) string {
	key := o.Type() + strconv.Itoa(o.Id())
	if _, exists := sfc.xobjs[key]; !exists {
		sfc.xobjs[key] = o
	}

	return key
}

func (sfc *pdfSurface) alterMatrix(a, b, c, d, e, f float64) {

	sfc.endText()

	fmt.Fprintf(sfc.w, "%f %f %f %f %f %f cm\r\n",
		a, b, c, d, e, f)
}

func (sfc *pdfSurface) Close() {
	sfc.endText()
}

func (sfc *pdfSurface) PushState() {

	sfc.endText()

	fmt.Fprint(sfc.w, "q\r\n")
}

func (sfc *pdfSurface) PopState() {
	sfc.endText()

	fmt.Fprint(sfc.w, "Q\r\n")
}

func (sfc *pdfSurface) Rotate(byRadians float64) {

	c := math.Cos(byRadians)
	s := math.Sin(byRadians)

	sfc.alterMatrix(c, s, -s, c, 0, 0)
}

func (sfc *pdfSurface) Translate(x, y float64) {

	sfc.alterMatrix(1, 0, 0, 1,
		d2g.ConvertUnit(x, sfc.u, d2g.U_PT),
		d2g.ConvertUnit(y, sfc.u, d2g.U_PT))
}

func (sfc *pdfSurface) Skew(xRadians float64, yRadians float64) {

	x := math.Atan(xRadians)
	y := math.Atan(yRadians)

	sfc.alterMatrix(1, x, y, 1, 0, 0)
}

func (sfc *pdfSurface) Scale(xScale float64, yScale float64) {

	sfc.alterMatrix(xScale, 0, 0, yScale, 0, 0)
}

func (sfc *pdfSurface) Fg(color d2g.Color) {

	sfc.endText()

	sfc.writeColor(color)
	fmt.Fprint(sfc.w, " RG\r\n")
}

func (sfc *pdfSurface) Bg(color d2g.Color) {

	sfc.endText()

	sfc.writeColor(color)
	fmt.Fprint(sfc.w, " rg\r\n")
}

func (sfc *pdfSurface) LineWidth(width float64) {
	sfc.endText()

	fmt.Fprintf(sfc.w, "%f w\r\n", d2g.ConvertUnit(width, sfc.u, d2g.U_PT))
}

func (sfc *pdfSurface) LineCap(capStyle d2g.LineCapStyle) {
	sfc.endText()

	fmt.Fprintf(sfc.w, "%d J\r\n", int32(capStyle))
}

func (sfc *pdfSurface) LineJoin(joinStyle d2g.LineJoinStyle) {
	sfc.endText()

	fmt.Fprintf(sfc.w, "%d j\r\n", int32(joinStyle))
}

func (sfc *pdfSurface) LinePattern(pattern []float64, phase float64) {
	sfc.endText()
}

func (sfc *pdfSurface) Stroke(path *d2g.Path) {
	sfc.endText()

	sfc.writePath(path)

	fmt.Fprint(sfc.w, "S\r\n")
}

func (sfc *pdfSurface) Fill(path *d2g.Path) {
	sfc.endText()

	sfc.writePath(path)

	fmt.Fprint(sfc.w, "f\r\n")
}

var charsToEscape = [8]rune{
	'\n', '\r', '\t', '\b', '\f', '(', ')', '\\',
}

var escapedChars = [8]byte{
	'n', 'r', 't', 'b', 'f', '(', ')', '\\',
}

const escapeChar = byte('\\')

func (sfc *pdfSurface) Text(f d2g.Font, x, y float64, text string) {

	if pf, ok := f.(*pdfFont); ok {

		if !sfc.inText {
			fmt.Fprint(sfc.w, "BT\r\n")
			sfc.inText = true
		}

		if sfc.lastFont == nil ||
			!sfc.lastFont.Equals(pf) {
			fmt.Fprintf(sfc.w, "/F%d %f Tf\r\n",
				pf.face.id,
				d2g.ConvertUnit(pf.size, sfc.u, d2g.U_PT))

			sfc.lastFont = pf
			sfc.fonts = append(sfc.fonts, pf.face)
		}

		fmt.Fprintf(sfc.w, "%f %f Td\r\n(",
			d2g.ConvertUnit(x, sfc.u, d2g.U_PT),
			d2g.ConvertUnit(y, sfc.u, d2g.U_PT))

		textBuf := new(bytes.Buffer)

		for _, c := range text {
			escaped := false
			for eIx, ec := range charsToEscape {
				if c == ec {
					textBuf.WriteByte(escapeChar)
					textBuf.WriteByte(escapedChars[eIx])
					escaped = true
				}
			}

			if !escaped {
				textBuf.WriteRune(c)
			}
		}
		if textBuf.Len() > 0 {
			sfc.w.Write(textBuf.Bytes())
		}

		fmt.Fprint(sfc.w, ") Tj\r\n")
	}
}

func (sfc *pdfSurface) Image(i d2g.Image, x, y, w, h float64) {

	sfc.endText()

	if pi, ok := i.(*pdfImage); ok {

		sfc.PushState()
		sfc.Translate(x, y)
		sfc.Scale(w, h)

		name := sfc.addXObj(pi)

		fmt.Fprintf(sfc.w, "/%s Do\r\n", name)

		sfc.PopState()
	}
}

func (sfc *pdfSurface) endText() {

	if sfc.inText {
		fmt.Fprint(sfc.w, "ET\r\n")
		sfc.inText = false
	}
}

func (sfc *pdfSurface) writeColor(c d2g.Color) {

	fmt.Fprintf(sfc.w, "%f %f %f",
		float64(c.R)/255.0,
		float64(c.G)/255.0,
		float64(c.B)/255.0)
}

func (sfc *pdfSurface) writePath(path *d2g.Path) {

	r := path.Reader()

	cmd, ok := r.ReadCommandType()
	for ok {

		switch cmd {

		case d2g.MoveCmdType:
			sfc.writeMove(
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT),
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT))

		case d2g.LineCmdType:
			sfc.writeLine(
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT),
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT))

		case d2g.CurveCmdType:
			sfc.writeCurve(
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT),
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT),
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT),
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT),
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT),
				d2g.ConvertUnit(r.ReadFloat64(), sfc.u, d2g.U_PT))

		case d2g.RectCmdType:
			x1 := r.ReadFloat64()
			y1 := r.ReadFloat64()
			fmt.Fprintf(sfc.w, "%f %f %f %f re\r\n",
				d2g.ConvertUnit(x1, sfc.u, d2g.U_PT),
				d2g.ConvertUnit(y1, sfc.u, d2g.U_PT),
				d2g.ConvertUnit(r.ReadFloat64()-x1, sfc.u, d2g.U_PT),
				d2g.ConvertUnit(r.ReadFloat64()-y1, sfc.u, d2g.U_PT))

		case d2g.ArcCmdType:
			sfc.writeArc(r)

		case d2g.CloseCmdType:
			fmt.Fprint(sfc.w, "h\r\n")

		default:
			r.Dump()
			panic("Unknown path command.")
		}

		cmd, ok = r.ReadCommandType()
	}
}

func (sfc *pdfSurface) writeMove(x float64, y float64) {
	fmt.Fprintf(sfc.w, "%f %f m\r\n", x, y)
}

func (sfc *pdfSurface) writeLine(x float64, y float64) {
	fmt.Fprintf(sfc.w, "%f %f l\r\n", x, y)
}

func (sfc *pdfSurface) writeCurve(x1, y1, x2, y2, x3, y3 float64) {
	fmt.Fprintf(sfc.w, "%f %f %f %f %f %f c\r\n",
		x1, y1, x2, y2, x3, y3)
}

/*
// Approximate a circular arc using multiple
// cubic Bézier curves, one for each π/2 segment.
//
// This is from:
//      http://hansmuller-flex.blogspot.com/2011/04/approximating-circular-arc-with-cubic.html
func arc(p *pdf.Path, comp vg.PathComp) {
        x0 := comp.X + comp.Radius*vg.Length(math.Cos(comp.Start))
        y0 := comp.Y + comp.Radius*vg.Length(math.Sin(comp.Start))
        p.Line(pdfPoint(x0, y0))

        a1 := comp.Start
        end := a1 + comp.Angle
        sign := 1.0
        if end < a1 {
                sign = -1.0
        }
        left := math.Abs(comp.Angle)
        for left > 0 {
                a2 := a1 + sign*math.Min(math.Pi/2, left)
                partialArc(p, comp.X, comp.Y, comp.Radius, a1, a2)
                left -= math.Abs(a2 - a1)
                a1 = a2
        }
}*/

func (sfc *pdfSurface) writeArc(rdr d2g.PathReader) {

	x := d2g.ConvertUnit(rdr.ReadFloat64(), sfc.u, d2g.U_PT)
	y := d2g.ConvertUnit(rdr.ReadFloat64(), sfc.u, d2g.U_PT)
	r := d2g.ConvertUnit(rdr.ReadFloat64(), sfc.u, d2g.U_PT)
	start := rdr.ReadFloat64()
	sweep := rdr.ReadFloat64()

	x0 := x + r*math.Cos(start)
	y0 := y + r*math.Sin(start)

	sfc.writeMove(x0, y0)

	a1 := start
	end := start + sweep
	sign := 1.0
	if end < a1 {
		sign = -1.0
	}
	left := math.Abs(sweep)
	for left > 0 {
		a2 := a1 + sign*math.Min(math.Pi/2, left)
		sfc.writePartialArc(x, y, r, a1, a2)
		left -= math.Abs(a2 - a1)
		a1 = a2
	}
}

func (sfc *pdfSurface) writePartialArc(x, y, r, a1, a2 float64) {

	a := (a2 - a1) / 2
	x4 := r * math.Cos(a)
	y4 := r * math.Sin(a)
	x1 := x4
	y1 := -y4

	const k = 0.5522847498 // some magic constant
	f := k * math.Tan(a)
	x2 := x1 + f*y4
	y2 := y1 + f*x4
	x3 := x2
	y3 := -y2

	// Rotate and translate points into position.
	ar := a + a1
	sinar := math.Sin(ar)
	cosar := math.Cos(ar)
	x2r := x2*cosar - y2*sinar + x
	y2r := x2*sinar + y2*cosar + y
	x3r := x3*cosar - y3*sinar + x
	y3r := x3*sinar + y3*cosar + y
	x4 = r*math.Cos(a2) + x
	y4 = r*math.Sin(a2) + y
	sfc.writeCurve(x2r, y2r, x3r, y3r, x4, y4)
}

/*

// Approximate a circular arc of fewer than π/2
// radians with cubic Bézier curve.
func partialArc(p *pdf.Path, x, y, r vg.Length, a1, a2 float64) {
        a := (a2 - a1) / 2
        x4 := r * vg.Length(math.Cos(a))
        y4 := r * vg.Length(math.Sin(a))
        x1 := x4
        y1 := -y4

        const k = 0.5522847498 // some magic constant
        f := k * vg.Length(math.Tan(a))
        x2 := x1 + f*y4
        y2 := y1 + f*x4
        x3 := x2
        y3 := -y2

        // Rotate and translate points into position.
        ar := a + a1
        sinar := vg.Length(math.Sin(ar))
        cosar := vg.Length(math.Cos(ar))
        x2r := x2*cosar - y2*sinar + x
        y2r := x2*sinar + y2*cosar + y
        x3r := x3*cosar - y3*sinar + x
        y3r := x3*sinar + y3*cosar + y
        x4 = r*vg.Length(math.Cos(a2)) + x
        y4 = r*vg.Length(math.Sin(a2)) + y
        p.Curve(pdfPoint(x2r, y2r), pdfPoint(x3r, y3r), pdfPoint(x4, y4))
}


*/
