package graph

import "fmt"

// Loss computes a loss function for the graph.
// Arranging will mean minimizing this loss value.
// l2 is the regularization weight, dist is the target distance for connected nodes.
// DistMin is a minimum distance for NON connected nodes.
func (g *Graph) Loss(l2, dist, distMin float64) (loss float64) {

	// Ensure  sizes remain reasonable (explicit l2 regularization)
	for _, x := range g.x {
		loss += x * x
	}
	for _, y := range g.y {
		loss += y * y
	}
	loss = loss * l2

	// Specific distance
	for i := 0; i < g.Size(); i++ {
		for j := i + 1; j < g.Size(); j++ {
			if g.Linked(i, j) {
				// Connected links
				d := g.Dist2(i, j)
				loss += (d - dist) * (d - dist)
			} else {
				// Not connected
				d := g.Dist2(i, j)
				if d < distMin {
					loss += distMin - d
				}
			}
		}
	}
	return loss
}

// DLoss computes the gradient of the loss function w.r.t. the x and y coordinates.
func (g *Graph) DLoss(l2, dist, distMin float64) (dx, dy []float64) {

	dx, dy = make([]float64, g.Size()), make([]float64, g.Size())
	for i, x := range g.x {
		dx[i] = l2 * 2 * x
	}
	for j, y := range g.y {
		dy[j] = l2 * 2 * y
	}

	for i := 0; i < g.Size(); i++ {
		for j := i + 1; j < g.Size(); j++ {
			if g.Linked(i, j) {
				d := g.Dist2(i, j)
				dx[i] += 2 * (d - dist) * (2 * g.x[i])
				dx[j] += 2 * (d - dist) * (2 * g.x[j])
				dy[i] += 2 * (d - dist) * (2 * g.y[i])
				dy[j] += 2 * (d - dist) * (2 * g.y[j])
			} else {
				d := g.Dist2(i, j)
				if d < distMin {
					dx[i] += -(2 * g.x[i])
					dx[j] += -(2 * g.x[j])
					dy[i] += -(2 * g.y[i])
					dy[j] += -(2 * g.y[j])
				}
			}
		}
	}
	return dx, dy
}

// Minimize will adjust the node position to minimize the loss function.
// lambda is the step, iter is the nbr of iterations.
func (g *Graph) Minimize(lambda, l2, dist, distMin float64, iter int) {
	for it := 0; it < iter; it++ {
		dx, dy := g.DLoss(l2, dist, distMin)
		for i := 0; i < g.Size(); i++ {
			g.x[i] += -lambda * dx[i]
			g.y[i] += -lambda * dy[i]
		}

		// Debug
		fmt.Printf("\n%d : loss = %f", it, g.Loss(l2, dist, distMin))
	}

}
