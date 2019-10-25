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
	// Size provides the nb of nodes of the graph
	Size() int
}

// EditGrapher is a graph that can be edited.
type EditGrapher interface {
	Grapher
	// Add a node, returning its id.
	Add(x, y float64, legend string) int
	// Remove a node from the graph.
	// All links to the removed node are also removed.
	Link(node1, node2 int)
	// Move node n to the new position.
	Move(node int, x, y float64)
}

// Arranger provides the detailled loss function and corresponding gradient.
type Arranger interface {
	// Arrange modifies the provided graph and retuns the resulting loss value.
	Arrange(g EditGrapher) float64
}
