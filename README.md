dox2go`
======

Dox2Go is a Go library for generating documents in a variety of formats. 

Currently PDF is the only supported output format.

Some Features
-------------

Path Drawing

* Lines
* Bezier Curves
* Rectangles
* Arcs
* Different line styles, joins and caps.

Transformations

* Rotation
* Skew
* Translate
* Scale

Image drawing, including transparency.

Example
-------

```go
package main

import (
	"bytes"
	"dox2go"
	"dox2go/pdf"
	"fmt"
	"os"
)

func main() {

	// Write the document to this buffer.
	var b bytes.Buffer

	// Create a Pdf Document instance.
	doc := pdf.NewPdfDoc(&b)

	// Add a page
	pWidth, pHeight := dox2go.StandardSize(dox2go.PS_A4, dox2go.U_MM) // A4 Page
	page := doc.CreatePage(dox2go.U_MM, // Work in millimeters
	pWidth, pHeight,
	dox2go.PO_Portrait)                             // Portrait orientation

	// Get the drawing surface of the page.
	s := page.Surface()

	// The document object manages font instances.
	font := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 20)
	s.Bg(dox2go.RGB(0, 0, 0)) // Text is drawn using the background colour.
	s.Text(font, 50, 100, "Hello World!")

	// We're finished so flush the document to the buffer.
	doc.Close()

	// Write it to a file.
	f, err := os.Create("tmp.pdf")
	if err != nil {
		fmt.Println("Could not create tmp.pdf")
		return
	}
	n, err := b.WriteTo(f)
	fmt.Printf("Wrote %d bytes\n", n)
}
```
