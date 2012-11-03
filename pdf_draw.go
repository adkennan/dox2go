package dox2go

import (
	"fmt"
	"io"
	"math"
)

type pdfSurface struct {
	w        io.Writer
	u        Unit
	inText   bool
	fonts    []*pdfTypeFace
	lastFont *pdfFont
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

func (sfc *pdfSurface) Translate(byDistance Point) {

	sfc.alterMatrix(1, 0, 0, 1,
		ConvertUnit(byDistance.X, sfc.u, U_PT),
		ConvertUnit(byDistance.Y, sfc.u, U_PT))
}

func (sfc *pdfSurface) Skew(xRadians float64, yRadians float64) {

	x := math.Atan(xRadians)
	y := math.Atan(yRadians)

	sfc.alterMatrix(1, x, y, 1, 0, 0)
}

func (sfc *pdfSurface) Scale(xScale float64, yScale float64) {

	sfc.alterMatrix(xScale, 0, 0, yScale, 0, 0)
}

func (sfc *pdfSurface) Fg(color Color) {

	sfc.endText()

	sfc.writeColor(color)
	fmt.Fprint(sfc.w, " RG\r\n")
}

func (sfc *pdfSurface) Bg(color Color) {

	sfc.endText()

	sfc.writeColor(color)
	fmt.Fprint(sfc.w, " rg\r\n")
}

func (sfc *pdfSurface) LineWidth(width float64) {
	sfc.endText()

	fmt.Fprintf(sfc.w, "%f w\r\n", ConvertUnit(width, sfc.u, U_PT))
}

func (sfc *pdfSurface) LineCap(capStyle LineCapStyle) {
	sfc.endText()

	fmt.Fprintf(sfc.w, "%d J\r\n", int32(capStyle))
}

func (sfc *pdfSurface) LineJoin(joinStyle LineJoinStyle) {
	sfc.endText()

	fmt.Fprintf(sfc.w, "%d j\r\n", int32(joinStyle))
}

func (sfc *pdfSurface) LinePattern(pattern []float64, phase float64) {
	sfc.endText()
}

func (sfc *pdfSurface) Stroke(path *Path) {
	sfc.endText()

	sfc.writePath(path)

	fmt.Fprint(sfc.w, "S\r\n")
}

func (sfc *pdfSurface) Fill(path *Path) {
	sfc.endText()

	sfc.writePath(path)

	fmt.Fprint(sfc.w, "f\r\n")
}

func (sfc *pdfSurface) Text(f Font, p Point, text string) {

	if pf, ok := f.(*pdfFont); ok {

		if !sfc.inText {
			fmt.Fprint(sfc.w, "BT\r\n")
			sfc.inText = true
		}

		if sfc.lastFont == nil ||
			!sfc.lastFont.Equals(pf) {
			fmt.Fprintf(sfc.w, "/F%d %f Tf\r\n",
				pf.face.id,
				ConvertUnit(pf.size, sfc.u, U_PT))

			sfc.lastFont = pf
			sfc.fonts = append(sfc.fonts, pf.face)
		}

		fmt.Fprintf(sfc.w, "%f %f Td\r\n",
			ConvertUnit(p.X, sfc.u, U_PT),
			ConvertUnit(p.Y, sfc.u, U_PT))

		fmt.Fprintf(sfc.w, "(%s) Tj\r\n", text)
	}
}

func (sfc *pdfSurface) endText() {

	if sfc.inText {
		fmt.Fprint(sfc.w, "ET\r\n")
		sfc.inText = false
	}
}

func (sfc *pdfSurface) writeColor(c Color) {

	fmt.Fprintf(sfc.w, "%f %f %f",
		float64(c.R)/255.0,
		float64(c.G)/255.0,
		float64(c.B)/255.0)
}

func (sfc *pdfSurface) writePath(path *Path) {

	for ix := range path.elements {
		cmd := path.elements[ix]
		switch cmd.t {

		case moveCmdType:
			fmt.Fprintf(sfc.w, "%f %f m\r\n",
				ConvertUnit(cmd.p[0].X, sfc.u, U_PT),
				ConvertUnit(cmd.p[0].Y, sfc.u, U_PT))

		case lineCmdType:
			fmt.Fprintf(sfc.w, "%f %f l\r\n",
				ConvertUnit(cmd.p[0].X, sfc.u, U_PT),
				ConvertUnit(cmd.p[0].Y, sfc.u, U_PT))

		case curveCmdType:
			fmt.Fprintf(sfc.w, "%f %f %f %f %f %f c\r\n",
				ConvertUnit(cmd.p[0].X, sfc.u, U_PT),
				ConvertUnit(cmd.p[0].Y, sfc.u, U_PT),
				ConvertUnit(cmd.p[1].X, sfc.u, U_PT),
				ConvertUnit(cmd.p[1].Y, sfc.u, U_PT),
				ConvertUnit(cmd.p[2].X, sfc.u, U_PT),
				ConvertUnit(cmd.p[2].Y, sfc.u, U_PT))

		case rectCmdType:
			fmt.Fprintf(sfc.w, "%f %f %f %f re\r\n",
				ConvertUnit(cmd.p[0].X, sfc.u, U_PT),
				ConvertUnit(cmd.p[0].Y, sfc.u, U_PT),
				ConvertUnit(cmd.p[1].X-cmd.p[0].X, sfc.u, U_PT),
				ConvertUnit(cmd.p[1].Y-cmd.p[0].Y, sfc.u, U_PT))

		case closeCmdType:
			fmt.Fprint(sfc.w, "h\r\n")
		}
	}
}
