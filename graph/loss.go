package graph

import (
	"fmt"
	"math"
	"math/rand"
)

// LossParam defines the parameter setting used for optimization.
type LossParam struct {
	l2, distTargt, distTargtW, distMin, distMinW, clW float64
	iter                                              int
	lambda                                            float64
}

// Loss computes a loss function for the graph.
// Arranging will mean minimizing this loss value.
// l2 is the regularization weight,
// distTargt is the TARGET distance for connected nodes,
// distTargtW is the weight associated with TARGET distance,
// distMin is a MINIMUM distance to try to keep for NON connected nodes,
// distMinW is the associated weight,
// clW is the cumulative length weight for connected nodes
func (g *Graph) Loss(lp *LossParam) (loss float64) {

	// Ensure  sizes remain reasonable (explicit l2 regularization)
	for _, x := range g.x {
		loss += x * x
	}
	for _, y := range g.y {
		loss += y * y
	}
	loss = loss * lp.l2

	for i := 0; i < g.Size(); i++ {
		for j := i + 1; j < g.Size(); j++ {
			if g.Linked(i, j) {
				// Connected links as compared to target
				d := g.Dist2(i, j)
				loss += lp.distTargtW * (d - lp.distTargt) * (d - lp.distTargt)

				// Cumulative length penalty
				loss += lp.clW * math.Sqrt(d)

			} else {
				// Not connected penalty if below distMin
				d := g.Dist2(i, j)
				if d < lp.distMin {
					loss += lp.distMinW * (lp.distMin - d)
				}
			}
		}
	}
	return loss
}

// DLoss computes the gradient of the loss function w.r.t. the x and y coordinates.
func (g *Graph) DLoss(lp *LossParam) (dx, dy []float64) {

	dx, dy = make([]float64, g.Size()), make([]float64, g.Size())
	for i, x := range g.x {
		dx[i] = lp.l2 * 2 * x
	}
	for j, y := range g.y {
		dy[j] = lp.l2 * 2 * y
	}

	for i := 0; i < g.Size(); i++ {
		for j := i + 1; j < g.Size(); j++ {
			if g.Linked(i, j) {
				// target length
				d := g.Dist2(i, j)
				dx[i] += lp.distTargtW * 2 * (d - lp.distTargt) * (2 * g.x[i])
				dx[j] += lp.distTargtW * 2 * (d - lp.distTargt) * (2 * g.x[j])
				dy[i] += lp.distTargtW * 2 * (d - lp.distTargt) * (2 * g.y[i])
				dy[j] += lp.distTargtW * 2 * (d - lp.distTargt) * (2 * g.y[j])

				// Cumulatif length penanlty
				dd := math.Sqrt(d)
				dx[i] += -lp.clW * g.x[i] / dd
				dx[j] += -lp.clW * g.x[j] / dd
				dy[i] += -lp.clW * g.y[i] / dd
				dy[j] += -lp.clW * g.y[j] / dd

			} else {
				d := g.Dist2(i, j)
				if d < lp.distMin {
					dx[i] += -lp.distMinW * (2 * g.x[i])
					dx[j] += -lp.distMinW * (2 * g.x[j])
					dy[i] += -lp.distMinW * (2 * g.y[i])
					dy[j] += -lp.distMinW * (2 * g.y[j])
				}
			}
		}
	}
	return dx, dy
}

// Minimize will adjust the node position to minimize the loss function.
// lambda is the step precision, iter is the nbr of iterations.
// Annealing is implemented with random values to attempt to capture global minimum, not being stuck with local minimum.
// Results are currently highly sensitive to the LossParam parameters.
func (g *Graph) Minimize(lp *LossParam) *Graph {
	for it := 1; it <= lp.iter; it++ {
		dx, dy := g.DLoss(lp)

		// Annealing factor
		ann := float64(lp.iter) / float64(it)

		for i := 0; i < g.Size(); i++ {
			g.x[i] += lp.lambda * (dx[i] + ann*rand.Float64())
			g.y[i] += lp.lambda * (dy[i] + ann*rand.Float64())
		}

		// Debug
		if it%(lp.iter/10) == 0 || it < 5 {
			fmt.Printf("\n%d : loss = %f", it, g.Loss(lp))
		}
	}
	return g
}
