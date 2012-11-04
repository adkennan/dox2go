/*
* dox2go - A document generating library for go.
*
* Copyright 2011 Andrew Kennan. All rights reserved.
*
* Redistribution and use in source and binary forms, with or without modification, are
* permitted provided that the following conditions are met:
*
* 1. Redistributions of source code must retain the above copyright notice, this list of
* conditions and the following disclaimer.
*
* 2. Redistributions in binary form must reproduce the above copyright notice, this list
* of conditions and the following disclaimer in the documentation and/or other materials
* provided with the distribution.
*
* THIS SOFTWARE IS PROVIDED BY ANDREW KENNAN ''AS IS'' AND ANY EXPRESS OR IMPLIED
* WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
* FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL ANDREW KENNAN OR
* CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
* CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
* SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
* ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
* NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
* ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*
* The views and conclusions contained in the software and documentation are those of the
* authors and should not be interpreted as representing official policies, either expressed
* or implied, of Andrew Kennan.
 */
package dox2go

import (
	"bytes"
	"os"
	"testing"
)

func TestSimplePage(t *testing.T) {
	var b bytes.Buffer

	d := NewPdfDoc(&b)
	page := d.CreatePage(U_MM, StandardSize(PS_A4, U_MM), PO_Portrait)

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
