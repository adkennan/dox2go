/*
* dox2go - A document generating library for go.
*
* Copyright 2013 Andrew Kennan. All rights reserved.
*
 */
package dox2go

import "testing"

func TestPath(t *testing.T) {

	p := NewPath()

	p.Move(Point{1, 1})
	p.Line(Point{10, 10})
	p.Curve(Point{5, 5}, Point{15, 15}, Point{20, 10})
	p.Close()

	const expected = 4
	if len(p.Elements) != expected {
		t.Errorf("Expected %d, was %d.", expected, len(p.Elements))
	}
}
