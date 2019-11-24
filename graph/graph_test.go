package graph

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestNormalize(t *testing.T) {

	rand.Seed(time.Now().Unix())

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

	g.Shuffle()
	g.Normalize()

	for i := range g.x {
		if g.x[i] > 1 || g.x[i] < -1 || g.y[i] > 1 || g.y[i] < -1 {
			log.Fatal(g.ToString())
		}
	}

}
func TestArrange(t *testing.T) {

	var err error

	g := NewWithLossCombined(NewGraph())
	g.lambda = 0.00001
	g.l2 = 0.001
	g.distTargt = 0.3
	g.distTargtW = 1.
	g.distMin = 0.3
	g.distMinW = 1.
	g.clW = 5.
	g.iter = 500

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

	g.Add(0.8, 5, "5")
	g.Add(0.6, 8, "6")
	g.Add(0.5, 3, "7")
	g.Link(6, 7)
	g.Link(7, 5)
	g.Link(5, 6)

	g.Link(1, 6)
	g.Link(6, 4)
	g.Link(5, 3)
	g.Link(7, 1)
	g.Link(5, 0)

	err = ioutil.WriteFile("ex.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	g.Normalize()

	err = ioutil.WriteFile("ex_normalized.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	g.Minimize()
	g.Normalize()

	err = ioutil.WriteFile("ex_minimized_normalized.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	g.Shuffle()
	g.Minimize()
	g.Normalize()
	err = ioutil.WriteFile("ex_shuffled_minimized_normalized.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
