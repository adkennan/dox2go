/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */
package dox2go

import (
	"image"
)

type PageSize int32

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

type PageOrientation int32

const (
	PO_Landscape PageOrientation = iota
	PO_Portrait
)

type Document interface {
	CreatePage(pu Unit, size Point, po PageOrientation) Page

	CreateFont(name string, fs FontStyle, size float64) Font

	CreateImage(src image.Image) Image

	Close() error
}

type Page interface {
	Surface() Surface
}
