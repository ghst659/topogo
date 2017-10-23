package topogo

import (
	"log"
	"os"
)

// Enum constants for graph edge directions.
const (
	kPre = iota
	kSuc
	kMax
)

// Returns the opposite direction of the given direction.
func opposite(d int) int {
	return (d + 1) % kMax
}

type nodeSet map[string]bool
type edgePair [kMax]nodeSet

// A representation of a directed graph that keeps track
// of directional edges between nodes.  Vertices are
// identified by strings.
type DirectedGraph struct {
	nodes map[string]edgePair
}

// Creates a new DirectedGraph and returns a pointer to
// the newly created graph.
func NewGraph() (g *DirectedGraph) {
	g = new(DirectedGraph)
	g.nodes = make(map[string]edgePair)
	return
}

// Adds a new node to the DirectedGraph.  If the node
// is already in the graph, nothing is done.
func (g *DirectedGraph) AddNode(n string) {
	if _, found := g.nodes[n]; ! found {
		g.nodes[n] = edgePair{
			kPre: make(nodeSet),
			kSuc: make(nodeSet),
		}
	}
}

// Deletes a node from the DirectedGraph.  If the node is
// not in the graph, nothing is done.
func (g *DirectedGraph) DelNode(n string) {
	if _, found := g.nodes[n]; found {
		g.removeLinks(n, kPre)
		g.removeLinks(n, kSuc)
		delete(g.nodes, n)
	}
}

// Returns true if the node is found in the DirectedGraph.
func (g *DirectedGraph) HasNode(n string) (exist bool) {
	_, exist = g.nodes[n]
	return
}

// Returns a slice of all nodes in the DirectedGraph.
// The order of the nodes is implementation-dependent.
func (g *DirectedGraph) AllNodes() (result []string) {
	for id, _ := range g.nodes {
		result = append(result, id)
	}
	return
}

// Creates an edge within the DirectedGraph from
// node 'a' to node 'b'.  If the nodes do not exist
// in the graph, they are first created.
func (g *DirectedGraph) AddEdge(a string, b string) {
	g.AddNode(a)
	g.AddNode(b)
	g.nodes[a][kSuc][b] = true
	g.nodes[b][kPre][a] = true
}

// Deletes the given edge from the DirectedGraph.
// If either of the nodes is not in the graph, or
// the edge does not exist, nothing is done.
func (g *DirectedGraph) DelEdge(a string, b string) {
	if g.HasNode(a) && g.HasNode(b) {
		delete(g.nodes[a][kSuc], b)
		delete(g.nodes[b][kPre], a)
	}
}

// Returns a slice of nodes that are the immediate
// successors of the given node.
func (g *DirectedGraph) Successors(n string) []string {
	return g.neighbours(n, kSuc)
}

// Returns a slice of nodes that are the immediate
// predecessors of the given node.
func (g *DirectedGraph) Predecessors(n string) []string {
	return g.neighbours(n, kPre)
}

// Returns a slice of nodes that are downstream
// of the given node.
func (g *DirectedGraph) Downstreams(n string) []string {
	return g.traceNodes([]string{n}, kSuc)
}

// Returns a slice of nodes that are upstream
// of the given node.
func (g *DirectedGraph) Upstreams(n string) []string {
	return g.traceNodes([]string{n}, kPre)
}

// Returns a slice of nodes that span the 'i' and 't' nodes.
func (g *DirectedGraph) Subgraph(i []string, t []string) (sg []string) {
	downstreams := g.traceNodes(i, kSuc)
	upstreams := g.traceNodes(t, kPre)
	sd := l2s(downstreams)
	su := l2s(upstreams)
	for n, _ := range sd {
		if su[n] {
			sg = append(sg, n)
		}
	}
	return
}

func (g *DirectedGraph) neighbours(n string, d int) (result []string) {
	if g.HasNode(n) {
		for id, _ := range g.nodes[n][d] {
			result = append(result, id)
		}
	}
	return
}

func (g *DirectedGraph) traceNodes(nodes []string, d int) (result []string) {
	tally := make(nodeSet)
	queue := make(chan string, 256)
	defer close(queue)
	for _, n := range nodes {
		queue <- n
	}
	queueLoop: for {
		select {
		case current := <-queue:
			if ! tally[current] {
				tally[current] = true
				for _, child := range g.neighbours(current, d) {
					queue <- child
				}
			}
		default:
			break queueLoop
		}
	}
	for n := range tally {
		result = append(result, n)
	}
	return
}

func (g *DirectedGraph) removeLinks(n string, d int) {
	if _, found := g.nodes[n]; found {
		o := opposite(d)
		for v, _ := range g.nodes[n][d] {
			delete(g.nodes[v][o], n)
		}
	}
}

func l2s(l []string) (m nodeSet) {
	m = make(nodeSet)
	for _, v := range l {
		m[v] = true
	}
	return
}

func init() {
	log.SetOutput(os.Stderr)
}
