package graph

import (
	"fmt"
	"math"
	"math/rand"
)

// Loss computes a loss function for the graph.
// Arranging will mean minimizing this loss value.
// l2 is the regularization weight,
// distTargt is the TARGET distance for connected nodes,
// distTargtW is the weight associated with TARGET distance,
// distMin is a MINIMUM distance to try to keep for NON connected nodes,
// distMinW is the associated weight,
// clW is the cumulative length weight for connected nodes
func (g *Graph) Loss(l2, distTargt, distTargtW, distMin, distMinW, clW float64) (loss float64) {

	// Ensure  sizes remain reasonable (explicit l2 regularization)
	for _, x := range g.x {
		loss += x * x
	}
	for _, y := range g.y {
		loss += y * y
	}
	loss = loss * l2

	for i := 0; i < g.Size(); i++ {
		for j := i + 1; j < g.Size(); j++ {
			if g.Linked(i, j) {
				// Connected links as compared to target
				d := g.Dist2(i, j)
				loss += distTargtW * (d - distTargt) * (d - distTargt)

				// Cumulative length penalty
				loss += clW * math.Sqrt(d)

			} else {
				// Not connected penalty if below distMin
				d := g.Dist2(i, j)
				if d < distMin {
					loss += distMinW * (distMin - d)
				}
			}
		}
	}
	return loss
}

// DLoss computes the gradient of the loss function w.r.t. the x and y coordinates.
func (g *Graph) DLoss(l2, distTargt, distTargtW, distMin, distMinW, clW float64) (dx, dy []float64) {

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
				// target length
				d := g.Dist2(i, j)
				dx[i] += distTargtW * 2 * (d - distTargt) * (2 * g.x[i])
				dx[j] += distTargtW * 2 * (d - distTargt) * (2 * g.x[j])
				dy[i] += distTargtW * 2 * (d - distTargt) * (2 * g.y[i])
				dy[j] += distTargtW * 2 * (d - distTargt) * (2 * g.y[j])

				// Cumulatif length penanlty
				dd := math.Sqrt(d)
				dx[i] += -clW * g.x[i] / dd
				dx[j] += -clW * g.x[j] / dd
				dy[i] += -clW * g.y[i] / dd
				dy[j] += -clW * g.y[j] / dd

			} else {
				d := g.Dist2(i, j)
				if d < distMin {
					dx[i] += -distMinW * (2 * g.x[i])
					dx[j] += -distMinW * (2 * g.x[j])
					dy[i] += -distMinW * (2 * g.y[i])
					dy[j] += -distMinW * (2 * g.y[j])
				}
			}
		}
	}
	return dx, dy
}

// Minimize will adjust the node position to minimize the loss function.
// lambda is the step, iter is the nbr of iterations.
func (g *Graph) Minimize(lambda, l2, distTargt, distTargtW, distMin, distMinW, clW float64, iter int) *Graph {
	for it := 1; it <= iter; it++ {
		dx, dy := g.DLoss(l2, distTargt, distTargtW, distMin, distMinW, clW)

		// Annealing factor
		ann := float64(iter) / float64(it)

		for i := 0; i < g.Size(); i++ {
			g.x[i] += lambda * (dx[i] + ann*rand.Float64())
			g.y[i] += lambda * (dy[i] + ann*rand.Float64())
		}

		// Debug
		if it%(iter/10) == 0 || it < 5 {
			fmt.Printf("\n%d : loss = %f", it, g.Loss(l2, distTargt, distTargtW, distMin, distMinW, clW))
		}
	}
	return g
}
