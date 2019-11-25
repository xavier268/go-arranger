package graph

import (
	"fmt"
	"math"
	"math/rand"
)

// WithLossCombined is a Graph that can be optimized
// for the provided Loss function, combining various penalties.
type WithLossCombined struct {
	*Graph
	*LossParam
}

// LossParam are the parameters used for the Loss and Minimization functions.
// L2 is the regularization weight,
// DistTargt is the TARGET distance for connected nodes,
// DistTargtW is the weight associated with TARGET distance,
// DistMin is a MINIMUM distance to try to keep for NON connected nodes,
// DistMinW is the associated weight,
// Clw is the cumulative length weight for connected nodes
// Repw is the weight for the repulsion force between any nodes
// AnnW is the Annealing Weight
type LossParam struct {
	L2, DistTargt, DistTargtW, DistMin, DistMinW, Clw, Repw float64
	AnnW                                                    float64
	Iter                                                    int
	Lambda                                                  float64
}

// NewWithLossCombined constructor from a Graph object.
func NewWithLossCombined(gg *Graph) *WithLossCombined {
	g := new(WithLossCombined)
	g.Graph = gg
	g.LossParam = new(LossParam)

	// default parameters
	g.Lambda = 0.00001
	g.L2 = 0.
	g.DistTargt = 1.
	g.DistTargtW = 1.
	g.DistMin = 1.
	g.DistMinW = 0.
	g.Clw = 1.
	g.Repw = 1.
	g.Iter = 1
	g.AnnW = 0.

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
			// For all node pair, repulsion cost
			d := g.Dist2(i, j)
			loss += g.Repw / d

			if g.Linked(i, j) {
				// Connected links as compared to target
				loss += g.DistTargtW * (d - g.DistTargt) * (d - g.DistTargt)

				// Cumulative length penalty
				loss += g.Clw * math.Sqrt(d)

			} else {
				// Not connected penalty if below DistMin
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
			// repulsion cost
			d := g.Dist2(i, j)
			d2 := d * d
			dx[i] += -2 * g.Repw * (g.x[i] - g.x[j]) / d2
			dx[j] += -2 * g.Repw * (g.x[j] - g.x[i]) / d2
			dy[i] += -2 * g.Repw * (g.y[i] - g.y[j]) / d2
			dy[j] += -2 * g.Repw * (g.y[j] - g.y[i]) / d2

			if g.Linked(i, j) {
				// target length

				dx[i] += g.DistTargtW * 2 * (d - g.DistTargt) * (2 * (g.x[i] - g.x[j]))
				dx[j] += g.DistTargtW * 2 * (d - g.DistTargt) * (2 * (g.x[j] - g.x[i]))
				dy[i] += g.DistTargtW * 2 * (d - g.DistTargt) * (2 * (g.y[i] - g.y[j]))
				dy[j] += g.DistTargtW * 2 * (d - g.DistTargt) * (2 * (g.y[j] - g.y[i]))

				// Cumulatif length penalty
				dd := math.Sqrt(d)
				dx[i] += g.Clw * (g.x[i] - g.x[j]) / dd
				dx[j] += g.Clw * (g.x[j] - g.x[i]) / dd
				dy[i] += g.Clw * (g.y[i] - g.y[j]) / dd
				dy[j] += g.Clw * (g.y[j] - g.y[i]) / dd

			} else {
				if d < g.DistMin {
					dx[i] += -g.DistMinW * (2 * (g.x[i] - g.x[j]))
					dx[j] += -g.DistMinW * (2 * (g.x[j] - g.x[i]))
					dy[i] += -g.DistMinW * (2 * (g.y[i] - g.y[j]))
					dy[j] += -g.DistMinW * (2 * (g.y[j] - g.y[i]))
				}
			}
		}
	}
	return dx, dy
}

// Clone makes a deep copy
func (g *WithLossCombined) Clone() *WithLossCombined {
	gg := NewWithLossCombined(g.Graph.Clone())
	gg.L2, gg.Lambda, gg.DistMin, gg.DistTargt = g.L2, g.Lambda, g.DistMin, g.DistTargt
	gg.DistMinW, gg.DistTargtW = g.DistMinW, g.DistTargtW
	gg.Clw, gg.Repw, gg.Iter, gg.AnnW = g.Clw, g.Repw, g.Iter, g.AnnW
	return gg
}

// DLossEst estimates the gradient of the loss function w.r.t. the x and y coordinates.
// Mainly used for debugging & testing, not very efficient.
func (g *WithLossCombined) DLossEst() (dx, dy []float64) {

	// Step used to estimate the gradient
	var step float64
	step = 1e-6

	dx, dy = make([]float64, g.Size()), make([]float64, g.Size())
	var gg *WithLossCombined

	for i := range dx {
		gg = g.Clone()
		gg.x[i] += step
		dx[i] = (gg.Loss() - g.Loss()) / step

		gg = g.Clone()
		gg.y[i] += step
		dy[i] = (gg.Loss() - g.Loss()) / step
	}

	return dx, dy
}

// Minimize will adjust the node position to minimize the loss function.
// Lambda is the step precision, Iter is the nbr of iterations.
// Annealing is implemented with random values to attempt to capture global minimum, not being stuck with local minimum.
// Results are currently highly sensitive to the LossParam parameters.
func (g *WithLossCombined) Minimize() {
	g.minimize(false)
}

// minimize using either estimate or calculation for gradient.
func (g *WithLossCombined) minimize(useEstimate bool) {
	for it := 1; it <= g.Iter; it++ {
		var dx, dy []float64
		if useEstimate {
			dx, dy = g.DLossEst()
		} else {
			dx, dy = g.DLoss()
		}

		// Annealing factor
		ann := g.AnnW * float64(g.Iter) / float64(it)

		for i := 0; i < g.Size(); i++ {
			g.x[i] += -g.Lambda * (dx[i] + ann*rand.Float64())
			g.y[i] += -g.Lambda * (dy[i] + ann*rand.Float64())
		}

		// Debug
		if it%(g.Iter/10) == 0 || it < 5 {
			var g2 float64
			for i := range dx {
				g2 += dx[i]*dx[i] + dy[i]*dy[i]
			}

			fmt.Printf("\n%d : \tloss = %e", it, g.Loss())
			fmt.Printf("\tGrad = %e", math.Sqrt(g2))
		}
	}

}
