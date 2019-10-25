package simplegraph

import "fmt"

// ToSVG returns a SVG representation of g.
func (g *Graph) ToSVG() string {

	r := 5                        // radius
	w := 2                        // stroke width
	var mx, my float64 = 400, 200 // Scale width and scale height

	s := `
	<svg version="1.1"
	baseProfile="full"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	xmlns:ev="http://www.w3.org/2001/xml-events">
	`

	for i := range g.x {
		x, y := int(mx*g.x[i]), int(my*g.y[i])
		l := g.Legend(i)
		s += fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"%d\" stroke=\"red\" fill=\"transparent\" stroke-width=\"%d\"/>", x, y, r, w)
		s += fmt.Sprintf("<text x=\"%d\" y=\"%d\" >%s</text>", x+r, y+r, l)

		for j := range g.x {
			if i < j && g.Linked(i, j) {
				xx, yy := int(mx*g.x[j]), int(my*g.y[j])
				s += fmt.Sprintf("<line x1=\"%d\" y1=\"%d\" x2=\"%d\" y2=\"%d\" stroke=\"blue\" stroke-width=\"%d\"/>", x, y, xx, yy, w)
			}
		}
	}
	return s + "</svg>"
}
