// Package simplegraph implements a basic but updatable graph.
package simplegraph

import (
	"fmt"

	"github.com/xavier268/arranger/common"
)

// Graph implements the EditGrapher interface.
type Graph struct {
	x, y, size []float64                   // nodes
	legend     []string                    // node legends
	links      map[struct{ x, y int }]bool // edges
}

// Compiler checks interface contract ...
var _ common.EditGrapher = new(Graph)

// NewGraph creates a new, empty Graph.
func NewGraph() *Graph {
	g := new(Graph)
	g.links = make(map[struct{ x, y int }]bool)
	return g
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
	s := struct{ x, y int }{i, j}
	b, ok := g.links[s]
	if ok && b {
		return true
	}
	return false
}

// ToSVG returns a SVG representation of g.
func (g *Graph) ToSVG() string {
	return "to be implemented"
}

// ToString human readable format (for debugging).
func (g *Graph) ToString() string {
	return fmt.Sprint(*g)
}

// Add adds a node.
func (g *Graph) Add(x, y float64, legend string) int {
	g.x, g.y = append(g.x, x), append(g.y, y)
	g.legend = append(g.legend, legend)
	return len(g.x) - 1
}

// Clone the Graph (deep copy).
func (g *Graph) Clone() *Graph {
	gg := NewGraph()
	for i := range g.x {
		ii := gg.Add(g.x[i], g.y[i], g.legend[i])
		if ii != i {
			panic("Inconstitent internal structure while cloning ?!")
		}
		for j := range g.x {
			if g.Linked(i, j) {
				gg.Link(i, j)
			}
		}
	}
	return gg
}

// Link establish a link between nodes i and j.
func (g *Graph) Link(i, j int) {
	g.links[struct{ x, y int }{i, j}] = true
}
