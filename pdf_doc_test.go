package dox2go

import (
	"bytes"
	"os"
	"testing"
)

func TestSimplePage(t *testing.T) {
	var b bytes.Buffer

	d := NewPdfDoc(&b)
	page := d.CreatePage(U_MM, PS_A4, PO_Portrait)

	s := page.Surface()

	s.LineWidth(3)

	p := NewPath()
	p.Move(Point{40, 40})
	p.Line(Point{140, 40})
	p.Curve(Point{120, 20}, Point{60, 20}, Point{40, 40})
	p.Close()

	s.Bg(RGB(255, 255, 0))
	s.Fill(p)

	s.Fg(RGB(255, 0, 0))
	s.Stroke(p)

	font := d.CreateFont(FONT_Helvetica, FS_Bold, 20)

	s.Bg(RGB(0, 0, 0))
	s.PushState()
	s.Translate(Point{50, 150})
	s.Skew(1.0, 0.5)
	s.Text(font, Point{0, 0}, "Hello")
	s.PopState()

	font = d.CreateFont(FONT_Times, FS_Regular, 30)

	s.Bg(RGB(100, 100, 100))
	s.Text(font, Point{101, 129}, "world")

	s.Bg(RGB(0, 200, 255))
	s.Text(font, Point{100, 130}, "world")

	d.Close()

	f, err := os.Create("tmp.pdf")
	if err != nil {
		t.Error("Could not create tmp.pdf")
	}
	b.WriteTo(f)

}
