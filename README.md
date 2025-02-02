[![GoDoc](https://godoc.org/github.com/lggomez/intree?status.svg)](https://godoc.org/github.com/lggomez/intree) [![Go Report Card](https://goreportcard.com/badge/github.com/lggomez/intree)](https://goreportcard.com/report/github.com/lggomez/intree)

# INTree for Go

**Fork of https://github.com/geozelot/intree with test coverage and minor improvements**

Very fast static, flat **IN**terval **Tree** implementation for reverse range searches.

Highly efficient and with almost no memory footprint other than the stored ranges.

Further scientific reading about the adapted algorithm and comparisons between different approaches (in C/C++) can be found [here](https://github.com/lh3/cgranges).


# Behaviour

* INTree will build the tree once (**static; no updates after creation**)
* INTree returns indices to the initial boundaries array
* INTree currently supports finding all interleaving boundaries for a single `float64` value

# Usage

## API ([GoDoc](https://godoc.org/github.com/lggomez/intree))

### `type Bounds`

`Bounds{}` is the main interface expected by `NewINTree()`; requires `Limits()` method to access interval limits.

```go
type Bounds interface {
    Limits() (lower, upper float64)
}
```

### `type ValuedBounds`

`ValuedBounds{}` is the main interface expected by `NewINTreeV()`, acting as a wrapper for `Bounds`; Expects the `Value()` method for retrieving a value associated with the given boundaries

```go
type ValuedBounds interface {
    Bounds
    Value() interface{}
}
```

### `type INTree`

`INTree{}` is the main package object; holds Slice of reference indices and the respective interval limits.

```go
type INTree struct {
    // contains filtered or unexported fields
}
```

### `func NewINTree`

`NewINTree()` is the main initialization function; creates the tree from the given Slice of Bounds.

```go
func NewINTree(bounds []Bounds) *INTree
```

### `func NewINTreeV`

`NewINTreeV()` is the main initialization function; creates the tree from the given Slice of ValuedBounds.

```go
func NewINTreeV(bounds []ValuedBounds) *INTree
```

### `func (*INTree) Including`

`Including()` is the main entry point for bounds searches; traverses the tree and collects intervals that overlap with the given value.

```go
func (t *INTree) Including(val float64) []int
```

## Import
```go
import (
    "github.com/lggomez/intree"
)
```

## Examples

#### Simple `Bounds{}` interface implementation:

```go
// SimpleBounds is a simple Struct implicitly implementing the Bounds interface.
type SimpleBounds struct {
  Lower, Upper float64
}

// Limits accesses the interval limits.
func (sb *SimpleBounds) Limits() (float64, float64) {
  return sb.Lower, sb.Upper
}
```

#### Test Setup:

```go
package main

import (
    "github.com/lggomez/intree"
    "fmt"
)

// define simple Struct holding interval limits
type SimpleBounds struct {
  Lower, Upper float64
}

  // add method to access limits; implicitly implements INTree.Bounds interface
  func (sb *SimpleBounds) Limits() (float64, float64) {
    return sb.Lower, sb.Upper
  }

func main() {

  // create typed var
  var tree *intree.INTree
  
  // create example bounds
  inputBounds := []intree.Bounds{

    &SimpleBounds{Lower: 4.0, Upper: 6.0},   // match
    &SimpleBounds{Lower: 5.0, Upper: 7.0},
    &SimpleBounds{Lower: 4.0, Upper: 8.0},   // match
    &SimpleBounds{Lower: 1.0, Upper: 3.0},
    &SimpleBounds{Lower: 7.0, Upper: 9.0},
    &SimpleBounds{Lower: 3.0, Upper: 6.0},   // match
    &SimpleBounds{Lower: 2.0, Upper: 3.0},
    &SimpleBounds{Lower: 5.3, Upper: 7.9},
    &SimpleBounds{Lower: 3.2, Upper: 7.5},   // match
    &SimpleBounds{Lower: 4.4, Upper: 5.1},
    &SimpleBounds{Lower: 4.1, Upper: 4.9},   // match
    &SimpleBounds{Lower: 4.1, Upper: 4.9},   // match, same interval
    &SimpleBounds{Lower: 1.3, Upper: 3.1},
    &SimpleBounds{Lower: 7.9, Upper: 8.9},

  }

  // initialize new INTree and create tree from inputBounds
  tree = intree.NewINTree(inputBounds)

  // parse return Slice with indices referencing inputBounds
  for _, matchedIndex := range tree.Including(4.3) {

    // using INTree.Bounds interface method to access limits
    lowerLimit, upperLimit := inputBounds[matchedIndex].Limits()

    fmt.Printf("Match at inputBounds index %2d with range [%.1f, %.1f]\n", matchedIndex, lowerLimit, upperLimit)

    /*
      Match at inputBounds index 11 with range [4.1, 4.9]
      Match at inputBounds index 10 with range [4.1, 4.9]
      Match at inputBounds index  5 with range [3.0, 6.0]
      Match at inputBounds index  2 with range [4.0, 8.0]
      Match at inputBounds index  0 with range [4.0, 6.0]
      Match at inputBounds index  9 with range [3.2, 7.5]
    */

  }

}
```

#### For more use cases you can see the [tests](intree_test.go)

#### Try on [Go Playground](https://play.golang.org/p/RRDavPcgyhx).

____

##### Inspired by this great [KDTree implementation](https://github.com/mourner/kdbush) for JavaScript and adapted from this excellent [Go port](https://github.com/MadAppGang/kdbush).
