/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */
package pdf

import (
	"bytes"
	d2g "dox2go"
	"image/png"
	"os"
	"testing"
)

func TestSimplePage(t *testing.T) {
	var b bytes.Buffer

	var f, err = os.Open("../gologo.png")
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	d := NewPdfDoc(&b)

	pi := d.CreateImage(img)

	page := d.CreatePage(d2g.U_MM, d2g.StandardSize(d2g.PS_A4, d2g.U_MM), d2g.PO_Portrait)

	s := page.Surface()

	s.LineWidth(3)

	p := d2g.NewPath()
	p.Move(d2g.Point{40, 40})
	p.Line(d2g.Point{140, 40})
	p.Curve(d2g.Point{120, 20}, d2g.Point{60, 20}, d2g.Point{40, 40})
	p.Close()

	s.Bg(d2g.RGB(255, 255, 0))
	s.Fill(p)

	s.Fg(d2g.RGB(255, 0, 0))
	s.Stroke(p)

	font := d.CreateFont(FONT_Helvetica, d2g.FS_Bold, 20)

	s.Bg(d2g.RGB(0, 0, 0))
	s.PushState()
	s.Translate(d2g.Point{50, 150})
	s.Skew(1.0, 0.5)
	s.Text(font, d2g.Point{0, 0}, "(Hello\\)")
	s.PopState()

	font = d.CreateFont(FONT_Times, d2g.FS_Regular, 30)

	s.Bg(d2g.RGB(100, 100, 100))
	s.Text(font, d2g.Point{101, 129}, "wor Ƶ ld")

	s.Image(pi, d2g.Point{40, 180}, d2g.Size{153, 55})

	//s.Bg(RGB(0, 200, 255))
	//s.Text(font, Point{100, 130}, "world")

	d.Close()

	f, err = os.Create("tmp.pdf")
	if err != nil {
		t.Error("Could not create tmp.pdf")
	}
	b.WriteTo(f)
}
