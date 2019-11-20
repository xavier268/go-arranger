package graph

import "fmt"

// ToSVG returns a SVG representation of g.
// Display is intended for normalized grph, ie with value between -1 and +1.
func (g *Graph) ToSVG() string {

	r := 5                         // radius for nodes
	w := 2                         // stroke width
	m := 30                        // margin of display
	var mx, my float64 = 1200, 600 // Width and  height of display in pixels

	s := `
	<svg version="1.1"
	baseProfile="full"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	xmlns:ev="http://www.w3.org/2001/xml-events">
	`

	s += fmt.Sprintf("\n<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" stroke=\"grey\" fill=\"transparent\" />", m, m, int(mx), int(my))

	for i := range g.x {
		// Print nodes
		x, y := m+int(mx*(g.x[i]+1)/2), m+int(my*(g.y[i]+1)/2)
		l := g.Legend(i)
		s += fmt.Sprintf("\n<circle cx=\"%d\" cy=\"%d\" r=\"%d\" stroke=\"red\" fill=\"transparent\" stroke-width=\"%d\"/>", x, y, r, w)
		s += fmt.Sprintf("\n<text x=\"%d\" y=\"%d\" >%s</text>", x+r, y+r, l)

		// print links
		for j := range g.x {
			if i < j && g.Linked(i, j) {
				xx, yy := m+int(mx*g.x[j]), m+int(my*g.y[j])
				s += fmt.Sprintf("\n<line x1=\"%d\" y1=\"%d\" x2=\"%d\" y2=\"%d\" stroke=\"blue\" stroke-width=\"%d\"/>", x, y, xx, yy, w)
			}
		}
	}
	return s + "\n</svg>"
}
