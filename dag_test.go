package topogo

import (
	"testing"
)

func TestAddNodesAndEdges(t *testing.T) {
	g := NewGraph()
	g.AddEdge("x", "p")
	g.AddEdge("x", "q")
	g.AddEdge("r", "y")
	g.AddEdge("q", "r")
	all := g.AllNodes()
	if !same(all, []string{"p", "x", "q", "y", "r"}) {
		t.Error("Want 5 nodes, got", all)
	}
}

func TestDelNode(t *testing.T) {
	g := NewGraph()
	g.AddEdge("x", "y")
	g.AddEdge("y", "z")
	g.DelNode("y")
	if !same(g.AllNodes(), []string{"x", "z"}) {
		t.Error("Failed to delete y")
	}
	if !same(g.Successors("x"), []string{}) {
		t.Error("X still has edges")
	}
	if !same(g.Predecessors("z"), []string{}) {
		t.Error("Z still has edges")
	}
}

func TestDelEdge(t *testing.T) {
	g := NewGraph()
	g.AddEdge("x", "y")
	g.AddEdge("y", "z")
	g.DelEdge("y", "z")
	if !same(g.AllNodes(), []string{"x", "y", "z"}) {
		t.Error("Unintended delete y")
	}
	if !same(g.Successors("y"), []string{}) {
		t.Error("Y still linked to Z")
	}
	if !same(g.Predecessors("z"), []string{}) {
		t.Error("Z still linked to Y")
	}
}

func TestSuccPred(t *testing.T) {
	g := NewGraph()
	g.AddEdge("x", "p")
	g.AddEdge("x", "q")
	g.AddEdge("y", "r")
	g.AddEdge("q", "r")

	xsuc := g.Successors("x")
	if !same(xsuc, []string{"p", "q"}) {
		t.Error("Invalid x successors:", xsuc)
	}

	rpre := g.Predecessors("r")
	if !same(rpre, []string{"y", "q"}) {
		t.Error("Invalid r predecessors:", rpre)
	}
	
	rsuc := g.Successors("r")
	if !same(rsuc, []string{}) {
		t.Error("Invalid r successors:", rsuc)
	}
}

func TestUpDownStreams(t *testing.T) {
	g := NewGraph()
	g.AddEdge("a", "x")
	g.AddEdge("a", "y")
	g.AddEdge("x", "p")
	g.AddEdge("x", "q")
	g.AddEdge("y", "r")
	g.AddEdge("p", "z")
	g.AddEdge("r", "z")
	adown := g.Downstreams("a")
	if !same(adown, []string{"a", "x", "y", "z", "p", "q", "r"}) {
		t.Error("Failed tracing downstreams of a", adown)
	}
}

func TestSubgraph(t *testing.T) {
	g := NewGraph()
	for _, t := range []string{"a", "b", "c"} {
		for _, m := range []string{"p", "q", "r"} {
			g.AddEdge(t, m)
		}
	}
	for _, m := range []string{"p", "q", "r"} {
		for _, b := range []string{"x", "y", "z"} {
			g.AddEdge(m, b)
		}
	}
	sg := g.Subgraph([]string{"b"}, []string{"y"})
	if !same(sg, []string{"b", "p", "q", "r", "y"}) {
		t.Error("bad subgraph:", sg)
	}
	
}

func same(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sa := l2s(a)
	for _, mb := range b {
		if ! sa[mb] {
			return false
		}
	}
	return true
}
