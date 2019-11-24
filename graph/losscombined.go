package graph

import (
	"fmt"
	"math"
	"math/rand"
)

// WithLossCombined is a Graph that can be optimized
// for the provided Loss function, combining various penalties.
// l2 is the regularization weight,
// distTargt is the TARGET distance for connected nodes,
// distTargtW is the weight associated with TARGET distance,
// distMin is a MINIMUM distance to try to keep for NON connected nodes,
// distMinW is the associated weight,
// clW is the cumulative length weight for connected nodes
type WithLossCombined struct {
	*Graph
	l2, distTargt, distTargtW, distMin, distMinW, clW float64
	iter                                              int
	lambda                                            float64
}

// NewWithLossCombined constructor from a Graph object.
func NewWithLossCombined(gg *Graph) *WithLossCombined {
	g := new(WithLossCombined)
	g.Graph = gg

	// default parameters
	g.lambda = 0.00001
	g.l2 = 0.001
	g.distTargt = 0.3
	g.distTargtW = 1.
	g.distMin = 0.3
	g.distMinW = 1.
	g.clW = 5.
	g.iter = 500

	return g
}

// Loss computes a loss function for the graph.
// Arranging will mean minimizing this loss value.
func (g *WithLossCombined) Loss() (loss float64) {

	// Ensure  sizes remain reasonable (explicit l2 regularization)
	for _, x := range g.x {
		loss += x * x
	}
	for _, y := range g.y {
		loss += y * y
	}
	loss = loss * g.l2

	for i := 0; i < g.Size(); i++ {
		for j := i + 1; j < g.Size(); j++ {
			if g.Linked(i, j) {
				// Connected links as compared to target
				d := g.Dist2(i, j)
				loss += g.distTargtW * (d - g.distTargt) * (d - g.distTargt)

				// Cumulative length penalty
				loss += g.clW * math.Sqrt(d)

			} else {
				// Not connected penalty if below distMin
				d := g.Dist2(i, j)
				if d < g.distMin {
					loss += g.distMinW * (g.distMin - d)
				}
			}
		}
	}
	return loss
}

// DLoss computes the gradient of the loss function w.r.t. the x and y coordinates.
func (g *WithLossCombined) DLoss() (dx, dy []float64) {

	dx, dy = make([]float64, g.Size()), make([]float64, g.Size())
	for i, x := range g.x {
		dx[i] = g.l2 * 2 * x
	}
	for j, y := range g.y {
		dy[j] = g.l2 * 2 * y
	}

	for i := 0; i < g.Size(); i++ {
		for j := i + 1; j < g.Size(); j++ {
			if g.Linked(i, j) {
				// target length
				d := g.Dist2(i, j)
				dx[i] += g.distTargtW * 2 * (d - g.distTargt) * (2 * g.x[i])
				dx[j] += g.distTargtW * 2 * (d - g.distTargt) * (2 * g.x[j])
				dy[i] += g.distTargtW * 2 * (d - g.distTargt) * (2 * g.y[i])
				dy[j] += g.distTargtW * 2 * (d - g.distTargt) * (2 * g.y[j])

				// Cumulatif length penanlty
				dd := math.Sqrt(d)
				dx[i] += -g.clW * g.x[i] / dd
				dx[j] += -g.clW * g.x[j] / dd
				dy[i] += -g.clW * g.y[i] / dd
				dy[j] += -g.clW * g.y[j] / dd

			} else {
				d := g.Dist2(i, j)
				if d < g.distMin {
					dx[i] += -g.distMinW * (2 * g.x[i])
					dx[j] += -g.distMinW * (2 * g.x[j])
					dy[i] += -g.distMinW * (2 * g.y[i])
					dy[j] += -g.distMinW * (2 * g.y[j])
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
func (g *WithLossCombined) Minimize() {
	for it := 1; it <= g.iter; it++ {
		dx, dy := g.DLoss()

		// Annealing factor
		ann := float64(g.iter) / float64(it)

		for i := 0; i < g.Size(); i++ {
			g.x[i] += g.lambda * (dx[i] + ann*rand.Float64())
			g.y[i] += g.lambda * (dy[i] + ann*rand.Float64())
		}

		// Debug
		if it%(g.iter/10) == 0 || it < 5 {
			fmt.Printf("\n%d : loss = %f", it, g.Loss())
		}
	}

}
