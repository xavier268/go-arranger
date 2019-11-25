package graph

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestLargeGraph(t *testing.T) {
	var err error

	g := NewWithLossCombined(NewGraph())
	g.Iter = 10000
	g.Lambda = 1e-8
	g.L2 = 1e-3
	g.AnnW = 1e4
	g.Repw = 1e3
	g.DistTargt = 0.1
	g.DistTargtW = 10.

	for i := 0; i < 80; i++ {
		g.Add(0, 0, fmt.Sprintf("%d", i))
		for j := 1; j < i; j++ {
			if GCD(i, j) != 1 {
				g.Link(j, i)
			}
		}
	}
	g.Shuffle()

	g.Normalize()

	err = ioutil.WriteFile("ex_large_normalized.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("\n-----------------")
	g.Minimize()
	g.Normalize()

	err = ioutil.WriteFile("ex_large_minimized_normalized.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("\n-----------------")

}

// GCD calculates GCD iteratively using remainder.
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
