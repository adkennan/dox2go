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

type Unit int32

const (
	U_PT Unit = 0
	U_MM      = 1
	U_CM      = 2
	U_IN      = 3
)

//       PT            MM          CM            IN
var convMatrix = [4][4]float64{
	{1.0, 2.83464567, 28.3464567, 72.0}, // PT
	{0.352777778, 1.0, 10.0, 25.4},      // MM
	{0.0352777778, 0.1, 1.0, 2.54},      // CM
	{1.0 / 72.0, 0.0394, 0.394, 1.0},    // IN
}

func ConvertUnit(val float64, from Unit, to Unit) float64 {
	if from < 0 || int(from) >= len(convMatrix) || to < 0 || int(to) >= len(convMatrix) {
		panic("Invalid Unit")
	}

	if from == to {
		return val
	}

	return convMatrix[to][from] * val
}
