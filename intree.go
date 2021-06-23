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
// Changelog: 	* Improve package code readability, add comments
//				* Add ValuedBounds interface and compatibility builder

// Package intree provides a very fast, static, flat, augmented interval tree for reverse range searches.
package intree

import (
	"math"
	"math/rand"
)

// Bounds is the main interface expected by NewINTree(); requires Limits method to access interval limits.
type Bounds interface {
	Limits() (lower, upper float64)
}

// ValuedBounds is the main interface expected by NewINTreeV(), acting as a wrapper for Bounds.
// Expects the Value() method for retrieving a value associated with the given boundaries
type ValuedBounds interface {
	Bounds
	Value() interface{}
}

// INTree is the main package object;
// holds Slice of reference indices and the respective interval limits.
type INTree struct {
	indexes []int
	limits  []float64
}

// NewINTree is the main initialization function;
// creates the tree from the given Slice of Bounds.
func NewINTree(bounds []Bounds) *INTree {
	tree := INTree{}
	tree.buildTree(bounds)

	return &tree
}

// NewINTreeV is the main initialization function;
// creates the tree from the given Slice of ValuedBounds.
func NewINTreeV(bounds []ValuedBounds) *INTree {
	tree := INTree{}
	tree.buildTreeV(bounds)

	return &tree
}

// buildTree is the internal tree construction function;
// creates, sorts and augments nodes into Slices.
func (t *INTree) buildTree(bounds []Bounds) {
	t.indexes = make([]int, len(bounds))
	t.limits = make([]float64, 3*len(bounds))

	for i, v := range bounds {
		t.indexes[i] = i
		l, u := v.Limits()

		t.limits[3*i] = l
		t.limits[3*i+1] = u
		t.limits[3*i+2] = 0
	}

	sort(t.limits, t.indexes)
	augment(t.limits, t.indexes)
}

// buildTreeV is the internal tree construction function for ValuedBounds;
// creates, sorts and augments nodes into Slices.
func (t *INTree) buildTreeV(bounds []ValuedBounds) {
	t.indexes = make([]int, len(bounds))
	t.limits = make([]float64, 3*len(bounds))

	for i, v := range bounds {
		t.indexes[i] = i
		l, u := v.Limits()

		t.limits[3*i] = l
		t.limits[3*i+1] = u
		t.limits[3*i+2] = 0
	}

	sort(t.limits, t.indexes)
	augment(t.limits, t.indexes)
}

// Including is the main entry point for bounds searches;
// traverses the tree and collects intervals that overlap with the given value.
func (t *INTree) Including(val float64) []int {
	idxStock := []int{0, len(t.indexes) - 1}
	result := []int{}

	for len(idxStock) > 0 {
		// Retrieve right and left boundaries from index stock
		rBoundIdx := idxStock[len(idxStock)-1]
		idxStock = idxStock[:len(idxStock)-1]
		lBoundIdx := idxStock[len(idxStock)-1]
		idxStock = idxStock[:len(idxStock)-1]

		if lBoundIdx == rBoundIdx+1 {
			continue
		}

		centerIdx := int(math.Ceil(float64(lBoundIdx+rBoundIdx) / 2.0))
		lowerLimit := t.limits[3*centerIdx+2]

		if val <= lowerLimit {
			idxStock = append(idxStock, lBoundIdx, centerIdx-1)
		}

		l := t.limits[3*centerIdx]

		if l <= val {
			idxStock = append(idxStock, centerIdx+1, rBoundIdx)

			upperLimit := t.limits[3*centerIdx+1]

			if val <= upperLimit {
				result = append(result, t.indexes[centerIdx])
			}
		}
	}

	return result
}

// augment is an internal utility function, adding maximum value of all child nodes to the current node.
func augment(limits []float64, indexes []int) {
	if len(indexes) < 1 {
		return
	}

	max := 0.0

	for idx := range indexes {
		if limits[3*idx+1] > max {
			max = limits[3*idx+1]
		}
	}

	r := len(indexes) >> 1

	limits[3*r+2] = max

	augment(limits[:3*r], indexes[:r])
	augment(limits[3*r+3:], indexes[r+1:])
}

// sort is an internal utility function, sorting the tree by lowest limits using Random Pivot QuickSearch
func sort(limits []float64, indexes []int) {
	if len(indexes) < 2 {
		return
	}

	// Pick index bounds
	l, r := 0, len(indexes)-1

	// Pick pivot
	p := rand.Int() % len(indexes)

	// Perform in-place assignment of limits and indexes
	indexes[p], indexes[r] = indexes[r], indexes[p]
	limits[3*p], limits[3*p+1], limits[3*p+2], limits[3*r], limits[3*r+1], limits[3*r+2] = limits[3*r], limits[3*r+1], limits[3*r+2], limits[3*p], limits[3*p+1], limits[3*p+2]

	for i := range indexes {
		if limits[3*i] < limits[3*r] {
			// Perform in-place assignment of limits and indexes
			indexes[l], indexes[i] = indexes[i], indexes[l]
			limits[3*l], limits[3*l+1], limits[3*l+2], limits[3*i], limits[3*i+1], limits[3*i+2] = limits[3*i], limits[3*i+1], limits[3*i+2], limits[3*l], limits[3*l+1], limits[3*l+2]

			l++
		}
	}

	// Perform in-place assignment of limits and indexes
	indexes[l], indexes[r] = indexes[r], indexes[l]
	limits[3*l], limits[3*l+1], limits[3*l+2], limits[3*r], limits[3*r+1], limits[3*r+2] = limits[3*r], limits[3*r+1], limits[3*r+2], limits[3*l], limits[3*l+1], limits[3*l+2]

	// Tail recursive calls on branches
	sort(limits[:3*l], indexes[:l])
	sort(limits[3*l+3:], indexes[l+1:])
}
