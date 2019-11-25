package graph

import "testing"

func BenchmarkGradientCalculus(b *testing.B) {

	g := NewWithLossCombined(NewGraph())

	g.Clw = 1.
	g.Repw = 1.
	g.DistTargtW = 1.
	g.DistMinW = 1.

	g.Add(1.654, 1., "0")
	g.Add(0.654, 25., "1")
	g.Add(2., 1., "2")
	g.Add(2., 2., "3")
	g.Link(3, 0)
	g.Link(1, 2)

	g.Shuffle()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.DLoss()
	}
}

func BenchmarkGradientEstimates(b *testing.B) {

	g := NewWithLossCombined(NewGraph())

	g.Clw = 1.
	g.Repw = 1.
	g.DistTargtW = 1.
	g.DistMinW = 1.

	g.Add(1.654, 1., "0")
	g.Add(0.654, 25., "1")
	g.Add(2., 1., "2")
	g.Add(2., 2., "3")
	g.Link(3, 0)
	g.Link(1, 2)

	g.Shuffle()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.DLossEst()
	}
}
