// Package common Common interfaces and types.
// All coordinates are float64, nodes are indentified by ints.
package common

// Grapher defines the interface for a displayable graph.
type Grapher interface {
	// Coord of the node
	Coord(node int) (x, y float64)
	// Linked is true if nodes are connected
	Linked(node1, node2 int) bool
	// Legend is a text describing the node.
	// This may contain svg or html code.
	Legend(node int) string
	// ToSVG renders the graph as an svg string fragment.
	ToSVG() string
	// ToString provides a human readable format for debugging.
	ToString() string
}

// EditGrapher is a graph that can be edited.
type EditGrapher interface {
	Grapher
	// Add a node, returning its id.
	Add(x, y float64, legend string) int
	// Remove a node from the graph.
	// All links to the removed node are also removed.
	Link(node1, node2 int)
}

// Arranger provides arranging capabilities to the graph.
// It never modifies an existing graph.
type Arranger interface {
	// Arrange returns a modified arranged copy of the original graph and the resukting loss value.
	// The original grap is not modified.
	Arrange(g Grapher) (Grapher, float64)
}

// Optimizer provides the detailled loss function and corresponding gradient.
type Optimizer interface {
	// Loss function associated with a given graph.
	Loss(g Grapher) float64
	// Gradient of the Loss function wrt the x and y positions of the nodes.
	Gradient(g Grapher) (dx, dy []float64)
}
