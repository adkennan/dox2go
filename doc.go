/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */

// dox2go is a library for generating documents.
package dox2go

import (
	"image"
)

// PageSize defines one of the standard page sizes defined below.
type PageSize int32

// These are the standard page sizes dox2go recognizes.
const (
	PS_A0 PageSize = iota
	PS_A1
	PS_A2
	PS_A3
	PS_A4
	PS_A5
	PS_A6
	PS_A7
	PS_A8
	PS_A9
	PS_A10
	PS_Letter
	PS_Legal
	PS_JuniorLegal
	PS_LedgerTabloid
)

var standardSizes = [15][2]float64{
	{841, 1189}, // A0
	{594, 841},  // A1
	{420, 594},  // A2
	{297, 420},  // A3
	{210, 297},  // A4
	{148, 210},  // A5
	{105, 148},  // A6
	{74, 105},   // A7
	{52, 74},    // A8
	{37, 52},    // A9
	{26, 37},    // A10
	{216, 279},  // Letter
	{216, 356},  // Legal
	{127, 203},  // Junior Legal
	{279, 432},  // Ledger/Tabloid
}

// StandardSize returns a page size in the requested units.
func StandardSize(ps PageSize, unit Unit) Point {
	if ps < 0 || int(ps) >= len(standardSizes) {
		panic("Invalid Page Size")
	}

	var s = standardSizes[ps]

	if unit == U_MM {
		return Point{s[0], s[1]}
	}

	return Point{
		ConvertUnit(s[0], U_MM, unit),
		ConvertUnit(s[1], U_MM, unit),
	}
}

// PageOrientation defines whether pages will be 
// oriented in portrait or landscape.
type PageOrientation int32

// These are the available page orientations.
const (
	PO_Landscape PageOrientation = iota
	PO_Portrait
)

// Document is the core interface of dox2go. It is responsible 
// for managing reusable objects like images and fonts as well
// as generating pages when requested.
//
// CreatePage constructs a new page object and adds it to 
// the document.
//
// CreateFont generates a Font object based on the supplied
// name, style and size. The name should be one of the 
// standard PDF font names.
//
// CreateImage returns an object that can be used to draw
// a bitmap on a document. The returned Image can be used
// any number of times.
//
// Close is called when the document is complete and is
// ready to be written to an output target.
type Document interface {
	CreatePage(pu Unit, size Point, po PageOrientation) Page

	CreateFont(name string, fs FontStyle, size float64) Font

	CreateImage(src image.Image) Image

	Close() error
}

// Page is an interface that describes a page of a Document.
//
// Surface returns the drawing surface of the page.
// Surface defines operations for drawing graphics and
// text on the page.
type Page interface {
	Surface() Surface
}
