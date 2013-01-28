/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
* Simple hello world example.
 */
package main

import (
	"bytes"
	"dox2go"
	"dox2go/pdf"
	"fmt"
	"math"
	"os"
)

func main() {

	var b bytes.Buffer

	doc := pdf.NewPdfDoc(&b)

	w, h := dox2go.StandardSize(dox2go.PS_A4, dox2go.U_MM)

	page := doc.CreatePage(dox2go.U_MM, w, h, dox2go.PO_Portrait)

	s := page.Surface()

	s.LineWidth(2.0)

	p := dox2go.NewPath()
	p.Arc(110, 150, 25, 0, math.Pi*2)
	p.Close()
	s.Fg(dox2go.RGB(0, 0, 0))
	s.Stroke(p)
	s.Bg(dox2go.RGB(255, 255, 0))
	s.Fill(p)

	p = dox2go.NewPath()
	p.Move(99, 158)
	p.Arc(97, 158, 2, 0, math.Pi*2)
	p.Move(125, 158)
	p.Arc(123, 158, 2, 0, math.Pi*2)

	p.Move(97, 140)
	p.Curve(105, 130, 115, 130, 123, 140)
	p.Close()

	s.Bg(dox2go.RGB(255, 255, 255))
	s.Fill(p)
	s.Stroke(p)

	font := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 20)
	s.Bg(dox2go.RGB(0, 0, 0))
	s.Text(font, 50, 100, "Hello World!")

	doc.Close()

	f, err := os.Create("tmp.pdf")
	if err != nil {
		fmt.Println("Could not create tmp.pdf")
		return
	}
	n, err := b.WriteTo(f)
	fmt.Printf("Write %d bytes\n", n)
}
