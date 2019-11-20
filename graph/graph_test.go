package graph

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
)

func TestNormalize(t *testing.T) {

	rand.Seed(42)

	g := NewGraph()
	g.Add(0, 0, "")
	g.Add(0, 0, "")
	g.Add(0, 0, "")
	g.Add(0, 0, "")
	g.Add(0, 0, "")
	g.Add(0, 0, "")
	g.Add(0, 0, "")
	g.Add(0, 0, "")
	g.Add(0, 0, "")

	g.Shuffle().Normalize()

	for i := range g.x {
		if g.x[i] > 1 || g.x[i] < -1 || g.y[i] > 1 || g.y[i] < -1 {
			log.Fatal(g.ToString())
		}
	}

}
func TestArrange(t *testing.T) {

	var err error

	g := NewGraph()
	g.Add(1.654, 1., "0")
	g.Add(0.654, 25., "1")
	g.Add(2., 1., "2")
	g.Add(2., 2., "3")

	g.Link(3, 0)
	g.Link(1, 2)
	g.Link(0, 0) // Should have no effect

	if !g.Linked(0, 3) || !g.Linked(1, 2) || !g.Linked(3, 0) || !g.Linked(3, 0) {
		t.FailNow()
	}

	if g.Linked(0, 0) || g.Linked(1, 3) || g.Linked(3, 1) {
		t.Fail()
	}

	g.Link(2, 3)
	g.Link(0, 1)
	g.Add(0.5, 0.5, "4")
	g.Link(4, 0)
	g.Link(4, 1)
	g.Link(4, 2)
	g.Link(4, 3)

	err = ioutil.WriteFile("test.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	g.Normalize()

	err = ioutil.WriteFile("test_normalized.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	g.Minimize(0.01, 0.1, 0.1, 0.3, 2000).Normalize()

	err = ioutil.WriteFile("test_minimized_normalized.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
