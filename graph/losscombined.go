package graph

import (
	"fmt"
	"math"
	"math/rand"
)

// WithLossCombined is a Graph that can be optimized
// for the provided Loss function, combining various penalties.
// L2 is the regularization weight,
// DistTargt is the TARGET distance for connected nodes,
// DistTargtW is the weight associated with TARGET distance,
// DistMin is a MINIMUM distance to try to keep for NON connected nodes,
// DistMinW is the associated weight,
// Clw is the cumulative length weight for connected nodes
type WithLossCombined struct {
	*Graph
	L2, DistTargt, DistTargtW, DistMin, DistMinW, Clw float64
	Iter                                              int
	Lambda                                            float64
}

// NewWithLossCombined constructor from a Graph object.
func NewWithLossCombined(gg *Graph) *WithLossCombined {
	g := new(WithLossCombined)
	g.Graph = gg

	// default parameters
	g.Lambda = 0.00001
	g.L2 = 0.001
	g.DistTargt = 0.3
	g.DistTargtW = 1.
	g.DistMin = 0.3
	g.DistMinW = 1.
	g.Clw = 5.
	g.Iter = 500

	return g
}

// Loss computes a loss function for the graph.
// Arranging will mean minimizing this loss value.
func (g *WithLossCombined) Loss() (loss float64) {

	// Ensure  sizes remain reasonable (explicit L2 regularization)
	for _, x := range g.x {
		loss += x * x
	}
	for _, y := range g.y {
		loss += y * y
	}
	loss = loss * g.L2

	for i := 0; i < g.Size(); i++ {
		for j := i + 1; j < g.Size(); j++ {
			if g.Linked(i, j) {
				// Connected links as compared to target
				d := g.Dist2(i, j)
				loss += g.DistTargtW * (d - g.DistTargt) * (d - g.DistTargt)

				// Cumulative length penalty
				loss += g.Clw * math.Sqrt(d)

			} else {
				// Not connected penalty if below DistMin
				d := g.Dist2(i, j)
				if d < g.DistMin {
					loss += g.DistMinW * (g.DistMin - d)
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
		dx[i] = g.L2 * 2 * x
	}
	for j, y := range g.y {
		dy[j] = g.L2 * 2 * y
	}

	for i := 0; i < g.Size(); i++ {
		for j := i + 1; j < g.Size(); j++ {
			if g.Linked(i, j) {
				// target length
				d := g.Dist2(i, j)
				dx[i] += g.DistTargtW * 2 * (d - g.DistTargt) * (2 * g.x[i])
				dx[j] += g.DistTargtW * 2 * (d - g.DistTargt) * (2 * g.x[j])
				dy[i] += g.DistTargtW * 2 * (d - g.DistTargt) * (2 * g.y[i])
				dy[j] += g.DistTargtW * 2 * (d - g.DistTargt) * (2 * g.y[j])

				// Cumulatif length penanlty
				dd := math.Sqrt(d)
				dx[i] += -g.Clw * g.x[i] / dd
				dx[j] += -g.Clw * g.x[j] / dd
				dy[i] += -g.Clw * g.y[i] / dd
				dy[j] += -g.Clw * g.y[j] / dd

			} else {
				d := g.Dist2(i, j)
				if d < g.DistMin {
					dx[i] += -g.DistMinW * (2 * g.x[i])
					dx[j] += -g.DistMinW * (2 * g.x[j])
					dy[i] += -g.DistMinW * (2 * g.y[i])
					dy[j] += -g.DistMinW * (2 * g.y[j])
				}
			}
		}
	}
	return dx, dy
}

// Minimize will adjust the node position to minimize the loss function.
// Lambda is the step precision, Iter is the nbr of iterations.
// Annealing is implemented with random values to attempt to capture global minimum, not being stuck with local minimum.
// Results are currently highly sensitive to the LossParam parameters.
func (g *WithLossCombined) Minimize() {
	for it := 1; it <= g.Iter; it++ {
		dx, dy := g.DLoss()

		// Annealing factor
		ann := float64(g.Iter) / float64(it)

		for i := 0; i < g.Size(); i++ {
			g.x[i] += g.Lambda * (dx[i] + ann*rand.Float64())
			g.y[i] += g.Lambda * (dy[i] + ann*rand.Float64())
		}

		// Debug
		if it%(g.Iter/10) == 0 || it < 5 {
			fmt.Printf("\n%d : loss = %f", it, g.Loss())
		}
	}

}
