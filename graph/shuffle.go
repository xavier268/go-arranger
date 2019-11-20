package graph

import "math/rand"

// Shuffle will reshuffle the graph node positions.
// All values are normalized between -1 and 1.
func (g *Graph) Shuffle() *Graph {

	for i := range g.x {
		r := rand.Float64()
		g.x[i] = 2*r - 1
		r = rand.Float64()
		g.y[i] = 2*rand.Float64() - 1
	}
	return g
}
