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

import "testing"

func checkConv(t *testing.T, fromVal float64, from Unit, to Unit, expected float64) {
	v := ConvertUnit(fromVal, from, to)
	if v != expected {
		t.Errorf("Expected %f. Was %f", expected, v)
	}
}

func TestConvertUnit(t *testing.T) {

	checkConv(t, 1.0, U_CM, U_MM, 10.0)
	checkConv(t, 1.0, U_MM, U_CM, 0.1)
	checkConv(t, 1.0, U_PT, U_IN, 1.0/72.0)
	checkConv(t, 1.0, U_IN, U_PT, 72.0)

	checkConv(t, 5.0, U_CM, U_MM, 50.0)
	checkConv(t, 5.0, U_MM, U_CM, 0.5)
	checkConv(t, 5.0, U_PT, U_IN, (1.0/72.0)*5.0)
	checkConv(t, 5.0, U_IN, U_PT, 72.0*5.0)
}
