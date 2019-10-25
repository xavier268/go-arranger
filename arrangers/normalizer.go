package arrangers

import "github.com/xavier268/arranger/common"

// Normalizer is an arranger that will ensure all coordiantes are between 0 and 1.
type Normalizer struct{}

// Compiler checks ...
var _ common.Arranger = new(Normalizer)

// Arrange will normalize the provided graph.
func (n *Normalizer) Arrange(g common.EditGrapher) float64 {

	var minx, miny, maxx, maxy float64 // Min & max of current coord.
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

	for i := 0; i < g.Size(); i++ {
		x, y := g.Coord(i)
		xx, yy := x, y // to do !! Scale & shift !!
		g.Move(i, xx, yy)
	}

	return 0.
}
