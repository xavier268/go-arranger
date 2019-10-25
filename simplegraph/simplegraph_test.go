package simplegraph

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {

	g := NewGraph()
	g.Add(1., 1., "11")
	g.Add(1., 2., "12")
	g.Add(2., 1., "21")
	g.Add(2., 2., "22")

	g.Link(3, 0)
	g.Link(1, 2)
	g.Link(0, 0) // Should have no effect

	if !g.Linked(0, 3) || !g.Linked(1, 2) || !g.Linked(3, 0) || !g.Linked(3, 0) {
		t.FailNow()
	}

	if g.Linked(0, 0) || g.Linked(1, 3) || g.Linked(3, 1) {
		t.Fail()
	}

	fmt.Println(g.ToSVG())
	err := ioutil.WriteFile("test.svg", []byte(g.ToSVG()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

}
