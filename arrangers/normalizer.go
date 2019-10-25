package arrangers

import "github.com/xavier268/arranger/common"

// Normalizer is an arranger that will ensure all coordiantes are between 0 and 1.
type Normalizer struct{}

// NewNormalizer constructs a Normalizer.
func NewNormalizer() *Normalizer {
	return new(Normalizer)
}

// Compiler checks ...
var _ common.Arranger = new(Normalizer)

// Arrange will normalize the provided graph.
func (n *Normalizer) Arrange(g common.EditGrapher) float64 {

	// Min & max of current coord.
	minx, miny := g.Coord(0)
	maxx, maxy := minx, miny

	for i := 0; i < g.Size(); i++ {
		x, y := g.Coord(i)
		if x < minx {
			minx = x
		}
		if x > maxx {
			maxx = x
		}
		if y < miny {
			miny = y
		}
		if y > maxy {
			maxy = y
		}
	}

	sx, sy := maxx-minx, maxy-miny // scales

	for i := 0; i < g.Size(); i++ {
		x, y := g.Coord(i)
		xx, yy := x-minx, y-miny // shift
		xx, yy = xx/sx, yy/sy    // scale
		g.Move(i, xx, yy)
	}

	return 0.
}
