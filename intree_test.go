// MIT License
//
// Copyright (c) 2020 geozelot (André Siefken), 2021 Luis Gomez
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// Changelog: Add unit tests

// Package intree_test provides tests for the intree package.
package intree_test

import (
	"fmt"
	"testing"

	"github.com/geozelot/intree"
	"github.com/stretchr/testify/assert"
)

type testBounds struct {
	Lower, Upper float64
}

func (tb *testBounds) Limits() (float64, float64) {
	return tb.Lower, tb.Upper
}

// valuedBounds is a composite intree.Bounds element with a value getter
type valuedBounds interface {
	intree.Bounds
	Value() int
}

type valuedTestBounds struct {
	Lower, Upper float64
	value        int
}

func (wb *valuedTestBounds) Limits() (float64, float64) {
	return wb.Lower, wb.Upper
}

func (wb *valuedTestBounds) Value() int {
	return wb.value
}

func Test_Tree(t *testing.T) {
	t.Run("Case_Example", func(t *testing.T) {
		inputBounds := []intree.Bounds{
			&testBounds{Lower: 4.0, Upper: 6.0},
			&testBounds{Lower: 5.0, Upper: 7.0},
			&testBounds{Lower: 4.0, Upper: 8.0},
			&testBounds{Lower: 1.0, Upper: 3.0},
			&testBounds{Lower: 7.0, Upper: 9.0},
			&testBounds{Lower: 3.0, Upper: 6.0},
			&testBounds{Lower: 2.0, Upper: 3.0},
			&testBounds{Lower: 5.3, Upper: 7.9},
			&testBounds{Lower: 3.2, Upper: 7.5},
			&testBounds{Lower: 4.4, Upper: 5.1},
			&testBounds{Lower: 4.1, Upper: 4.9},
			&testBounds{Lower: 1.3, Upper: 3.1},
			&testBounds{Lower: 7.9, Upper: 8.9},
		}

		tree := intree.NewINTree(inputBounds)
		matches := tree.Including(4.3)

		assert.EqualValues(t, 5, len(matches))
		for _, matchedIndex := range matches {
			lowerLimit, upperLimit := inputBounds[matchedIndex].Limits()

			fmt.Printf("Match at inputBounds index %2d with range [%.1f, %.1f]\n", matchedIndex, lowerLimit, upperLimit)

			switch matchedIndex {
			case 0:
				assert.EqualValues(t, 4.0, lowerLimit)
				assert.EqualValues(t, 6.0, upperLimit)
			case 2:
				assert.EqualValues(t, 4.0, lowerLimit)
				assert.EqualValues(t, 8.0, upperLimit)
			case 5:
				assert.EqualValues(t, 3.0, lowerLimit)
				assert.EqualValues(t, 6.0, upperLimit)
			case 9:
				assert.EqualValues(t, 3.2, lowerLimit)
				assert.EqualValues(t, 7.5, upperLimit)
			case 10:
				assert.EqualValues(t, 4.1, lowerLimit)
				assert.EqualValues(t, 4.9, upperLimit)
			}
		}
	})
	t.Run("Case_Border/nil_bounds", func(t *testing.T) {
		tree := intree.NewINTree(nil)
		matches := tree.Including(4.3)
		assert.EqualValues(t, 0, len(matches))
	})
	t.Run("Case_Border/single_interval", func(t *testing.T) {
		inputBounds := []intree.Bounds{
			&testBounds{Lower: 4.0, Upper: 6.0},
		}

		tree := intree.NewINTree(inputBounds)

		matches := tree.Including(4.3)
		assert.EqualValues(t, 1, len(matches))

		matches = tree.Including(7)
		assert.EqualValues(t, 0, len(matches))
	})
	t.Run("Case_Border/repeated_interval", func(t *testing.T) {
		inputBounds := []intree.Bounds{
			&testBounds{Lower: 4.0, Upper: 6.0},
			&testBounds{Lower: 4.0, Upper: 6.0},
		}

		tree := intree.NewINTree(inputBounds)

		matches := tree.Including(4.3)
		assert.EqualValues(t, 2, len(matches))

		matches = tree.Including(7)
		assert.EqualValues(t, 0, len(matches))
	})
	t.Run("Case_Border/overlap_at_boundary", func(t *testing.T) {
		inputBounds := []intree.Bounds{
			&testBounds{Lower: 4.0, Upper: 6.0},
			&testBounds{Lower: 6.0, Upper: 9.0},
			&testBounds{Lower: 9.0, Upper: 11.0},
		}

		tree := intree.NewINTree(inputBounds)

		matches := tree.Including(6.0)
		assert.EqualValues(t, 2, len(matches))

		matches = tree.Including(9.0)
		assert.EqualValues(t, 2, len(matches))
	})
	t.Run("Case_Border/fine_grained", func(t *testing.T) {
		// As we are using IEEE 754 double precision floats, we will have pathological cases
		// in which rounding may cause spurious results on interval matching. Nevertheless
		// we perform some basic asserts with non-pathological, arbitrary precision numbers
		// to validate the library implementation
		inputBounds := []intree.Bounds{
			&testBounds{Lower: 3.43567981e-21, Upper: 3.43567984e-21},
			&testBounds{Lower: 3.43567987e-21, Upper: 3.43567990e-21},
		}

		tree := intree.NewINTree(inputBounds)

		matches := tree.Including(3.43567985e-21) // between intervals
		assert.EqualValues(t, 0, len(matches))

		matches = tree.Including(3.43567987e-21) // at boundary
		assert.EqualValues(t, 1, len(matches))
		matches = tree.Including(3.43567988e-21) // inside second interval
		assert.EqualValues(t, 1, len(matches))
	})
}

func Test_Tree_Valued(t *testing.T) {
	t.Run("Case_Example", func(t *testing.T) {
		inputBounds := []intree.Bounds{
			&valuedTestBounds{Lower: 4.0, Upper: 6.0, value: 1},
			&valuedTestBounds{Lower: 5.0, Upper: 7.0, value: 2},
			&valuedTestBounds{Lower: 4.0, Upper: 8.0, value: 3},
			&valuedTestBounds{Lower: 1.0, Upper: 3.0, value: 4},
			&valuedTestBounds{Lower: 7.0, Upper: 9.0, value: 5},
			&valuedTestBounds{Lower: 3.0, Upper: 6.0, value: 6},
			&valuedTestBounds{Lower: 2.0, Upper: 3.0, value: 7},
			&valuedTestBounds{Lower: 5.3, Upper: 7.9, value: 8},
			&valuedTestBounds{Lower: 3.2, Upper: 7.5, value: 9},
			&valuedTestBounds{Lower: 4.4, Upper: 5.1, value: 10},
			&valuedTestBounds{Lower: 4.1, Upper: 4.9, value: 11},
			&valuedTestBounds{Lower: 1.3, Upper: 3.1, value: 12},
			&valuedTestBounds{Lower: 7.9, Upper: 8.9, value: 13},
		}

		tree := intree.NewINTree(inputBounds)
		matches := tree.Including(4.3)

		assert.EqualValues(t, 5, len(matches))
		for _, matchedIndex := range matches {
			lowerLimit, upperLimit := inputBounds[matchedIndex].Limits()

			fmt.Printf("Match at inputBounds index %2d with range [%.1f, %.1f]\n", matchedIndex, lowerLimit, upperLimit)

			switch matchedIndex {
			case 0:
				assert.EqualValues(t, 4.0, lowerLimit)
				assert.EqualValues(t, 6.0, upperLimit)
				assert.EqualValues(t, matchedIndex+1, inputBounds[matchedIndex].(valuedBounds).Value())
			case 2:
				assert.EqualValues(t, 4.0, lowerLimit)
				assert.EqualValues(t, 8.0, upperLimit)
				assert.EqualValues(t, matchedIndex+1, inputBounds[matchedIndex].(valuedBounds).Value())
			case 5:
				assert.EqualValues(t, 3.0, lowerLimit)
				assert.EqualValues(t, 6.0, upperLimit)
				assert.EqualValues(t, matchedIndex+1, inputBounds[matchedIndex].(valuedBounds).Value())
			case 9:
				assert.EqualValues(t, 3.2, lowerLimit)
				assert.EqualValues(t, 7.5, upperLimit)
				assert.EqualValues(t, matchedIndex+1, inputBounds[matchedIndex].(valuedBounds).Value())
			case 10:
				assert.EqualValues(t, 4.1, lowerLimit)
				assert.EqualValues(t, 4.9, upperLimit)
				assert.EqualValues(t, matchedIndex+1, inputBounds[matchedIndex].(valuedBounds).Value())
			}
		}
	})
	t.Run("Case_Border/nil_bounds", func(t *testing.T) {
		tree := intree.NewINTree(nil)
		matches := tree.Including(4.3)
		assert.EqualValues(t, 0, len(matches))
	})
	t.Run("Case_Border/single_interval", func(t *testing.T) {
		inputBounds := []intree.Bounds{
			&valuedTestBounds{Lower: 4.0, Upper: 6.0, value: 1},
		}

		tree := intree.NewINTree(inputBounds)

		matches := tree.Including(4.3)
		assert.EqualValues(t, 1, len(matches))
		assert.EqualValues(t, matches[0]+1, inputBounds[matches[0]].(valuedBounds).Value())

		matches = tree.Including(7)
		assert.EqualValues(t, 0, len(matches))
	})
	t.Run("Case_Border/repeated_interval", func(t *testing.T) {
		inputBounds := []intree.Bounds{
			&valuedTestBounds{Lower: 4.0, Upper: 6.0, value: 1},
			&valuedTestBounds{Lower: 4.0, Upper: 6.0, value: 2},
		}

		tree := intree.NewINTree(inputBounds)

		matches := tree.Including(4.3)
		assert.EqualValues(t, 2, len(matches))
		assert.EqualValues(t, matches[0]+1, inputBounds[matches[0]].(valuedBounds).Value())
		assert.EqualValues(t, matches[1]+1, inputBounds[matches[1]].(valuedBounds).Value())

		matches = tree.Including(7)
		assert.EqualValues(t, 0, len(matches))
	})
	t.Run("Case_Border/overlap_at_boundary", func(t *testing.T) {
		inputBounds := []intree.Bounds{
			&valuedTestBounds{Lower: 4.0, Upper: 6.0, value: 1},
			&valuedTestBounds{Lower: 6.0, Upper: 9.0, value: 2},
			&valuedTestBounds{Lower: 9.0, Upper: 11.0, value: 3},
		}

		tree := intree.NewINTree(inputBounds)

		matches := tree.Including(6.0)
		assert.EqualValues(t, 2, len(matches))
		assert.EqualValues(t, matches[0]+1, inputBounds[matches[0]].(valuedBounds).Value())
		assert.EqualValues(t, matches[1]+1, inputBounds[matches[1]].(valuedBounds).Value())

		matches = tree.Including(9.0)
		assert.EqualValues(t, 2, len(matches))
		assert.EqualValues(t, matches[0]+1, inputBounds[matches[0]].(valuedBounds).Value())
		assert.EqualValues(t, matches[1]+1, inputBounds[matches[1]].(valuedBounds).Value())
	})
	t.Run("Case_Border/fine_grained", func(t *testing.T) {
		// As we are using IEEE 754 double precision floats, we will have pathological cases
		// in which rounding may cause spurious results on interval matching. Nevertheless
		// we perform some basic asserts with non-pathological, arbitrary precision numbers
		// to validate the library implementation
		inputBounds := []intree.Bounds{
			&valuedTestBounds{Lower: 3.43567981e-21, Upper: 3.43567984e-21, value: 1},
			&valuedTestBounds{Lower: 3.43567987e-21, Upper: 3.43567990e-21, value: 2},
		}

		tree := intree.NewINTree(inputBounds)

		matches := tree.Including(3.43567985e-21) // between intervals
		assert.EqualValues(t, 0, len(matches))

		matches = tree.Including(3.43567987e-21) // at boundary
		assert.EqualValues(t, 1, len(matches))
		assert.EqualValues(t, matches[0]+1, inputBounds[matches[0]].(valuedBounds).Value())

		matches = tree.Including(3.43567988e-21) // inside second interval
		assert.EqualValues(t, 1, len(matches))
		assert.EqualValues(t, matches[0]+1, inputBounds[matches[0]].(valuedBounds).Value())
	})
}
