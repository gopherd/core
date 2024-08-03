package graphviz

import (
	"os"
	"testing"
)

func TestGraphviz(t *testing.T) {
	g := New("test", Directed)
	g.AddPlain("p1", "p2", "a")
	g.AddPlain("p1", "p3", "a")

	g.Add(NewEntity("n1", `[shape=box]`), NewEntity("n2", `[color=red, shape=box]`), "a")

	os.MkdirAll("./testdata", 0755)
	g.WriteFile("./testdata/out.gv")
}
