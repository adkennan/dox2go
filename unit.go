/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */

package dox2go

// Unit describes a unit of measurement.
type PageUnit int32

// Available measurement unit types.
const (
	U_PT PageUnit = 0
	U_MM          = 1
	U_CM          = 2
	U_IN          = 3
)

//       PT            MM          CM            IN
var convMatrix = [4][4]float64{
	{1.0, 2.83464567, 28.3464567, 72.0}, // PT
	{0.352777778, 1.0, 10.0, 25.4},      // MM
	{0.0352777778, 0.1, 1.0, 2.54},      // CM
	{1.0 / 72.0, 0.0394, 0.394, 1.0},    // IN
}

// ConvetUnit converts a value from one unit to another.
func ConvertUnit(val float64, from PageUnit, to PageUnit) float64 {
	if from < 0 || int(from) >= len(convMatrix) || to < 0 || int(to) >= len(convMatrix) {
		panic("Invalid Unit")
	}

	if from == to {
		return val
	}

	return convMatrix[to][from] * val
}
