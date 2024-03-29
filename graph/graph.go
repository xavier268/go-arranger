// Package graph implements a basic but updatable graph.
package graph

import "fmt"

// Graph implements a basic, non directed graph.
type Graph struct {
	x, y   []float64                   // nodes coords
	legend []string                    // node legends
	links  map[struct{ x, y int }]bool // edges, encoded with i<j
}

// NewGraph creates a new, empty Graph.
func NewGraph() *Graph {
	g := new(Graph)
	g.links = make(map[struct{ x, y int }]bool)
	return g
}

// Clone make a (deep) clone
func (g *Graph) Clone() *Graph {
	gg := NewGraph()
	for i := range g.x {
		gg.x, gg.y, gg.legend = append(gg.x, g.x[i]), append(gg.y, g.y[i]), append(gg.legend, g.legend[i])

		for j := range g.x {
			if g.Linked(i, j) {
				gg.Link(i, j)
			}
		}

	}

	return gg
}

// Size retun the number of nodes.
// Nodes are numbered from O to Size -1.
func (g *Graph) Size() int {
	return len(g.x)
}

// Coord gets the position of node n.
func (g *Graph) Coord(n int) (float64, float64) {
	return g.x[n], g.y[n]
}

// Legend provides the legend associated with the node.
func (g *Graph) Legend(n int) string {
	return g.legend[n]
}

// Linked is true  if nodes i and j are connected.
func (g *Graph) Linked(i, j int) bool {
	if i == j {
		return false
	}
	if i > j {
		return g.Linked(j, i)
	}
	s := struct{ x, y int }{i, j}
	b, ok := g.links[s]
	if ok && b {
		return true
	}
	return false
}

// ToString human readable format (for debugging).
func (g *Graph) ToString() string {
	return fmt.Sprint(*g)
}

// Add adds a node, returning its index.
func (g *Graph) Add(x, y float64, legend string) int {
	g.x, g.y = append(g.x, x), append(g.y, y)
	g.legend = append(g.legend, legend)
	return len(g.x) - 1
}

// Link establish a link between nodes i and j.
// Self linking is ignored.
// Links are symetrical, if (i,j) are linked, (j,i) will be.
func (g *Graph) Link(i, j int) {
	if i == j {
		return
	}
	if i > j {
		g.Link(j, i)
		return
	}
	g.links[struct{ x, y int }{i, j}] = true
}

// Dist2 provides the squared distance between two nodes
func (g *Graph) Dist2(i, j int) float64 {
	return (g.x[i]-g.x[j])*(g.x[i]-g.x[j]) + (g.y[i]-g.y[j])*(g.y[i]-g.y[j])
}
