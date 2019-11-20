package graph

// Normalize graph, bringing all nodes values between -1 and 1.
func (g *Graph) Normalize() *Graph {

	if g.Size() == 0 {
		return
	}

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
		// xx,yy are between 0 and 1.
		g.x[i], g.y[i] = 2*xx-1, 2*yy-1
	}
	return g
}
