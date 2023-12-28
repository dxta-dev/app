package graphs

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

type Node struct {
	Position Point
}

type Graph struct {
	Nodes []*Node
	Edges [][2]int
}

func Distance(n1, n2 Point) float64 {
	dx := n1.X - n2.X
	dy := n1.Y - n2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func NewNode(x, y float64) *Node {
	return &Node{Position: Point{X: x, Y: y}}
}

func NewGraph() *Graph {
	return &Graph{}
}

func (g *Graph) AddNode(node *Node) {
	g.Nodes = append(g.Nodes, node)
}

func (g *Graph) AddEdge(v1, v2 int) {
	g.Edges = append(g.Edges, [2]int{v1, v2})
}

func repulsiveForce(distance float64) float64 {
	if distance == 0 {
		distance = 0.1
	}
	k := 3600 * 3600 * 3600 * 600.0
	return k / (distance * distance * distance)
}

func ForceDirectedGraphLayout(graph *Graph, iterations int) {
	temperature := 1.0
	maxDisplacement := 60 * 60 * 2.0

	graph.createEdges()

	for iter := 0; iter < iterations; iter++ {
		displacement := make([]Point, len(graph.Nodes))

		for _, edge := range graph.Edges {
			v := graph.Nodes[edge[0]]
			u := graph.Nodes[edge[1]]
			delta := Point{X: v.Position.X - u.Position.X, Y: v.Position.Y - u.Position.Y}
			distance := math.Sqrt(delta.X*delta.X + delta.Y*delta.Y)
			force := repulsiveForce(distance)
			fmt.Println("force", edge[0], edge[1], distance, force)
			if delta.Y >= 0 {
				displacement[edge[0]].Y += force
				displacement[edge[1]].Y -= force
			} else {
				displacement[edge[0]].Y -= force
				displacement[edge[1]].Y += force
			}
		}

		for i, v := range graph.Nodes {
			dispLength := math.Abs(displacement[i].Y)
			if displacement[i].Y != 0 {
				fmt.Println("disp", i, displacement[i].Y / dispLength * dispLength * temperature)
				v.Position.Y += displacement[i].Y / dispLength * math.Min(dispLength * temperature, maxDisplacement)
			}
		}

		graph.removeEdges()
		graph.createEdges()

		temperature *= 0.95

	}
}

func (g *Graph) removeEdges() {
	for i := range g.Edges {
		g.Edges[i] = [2]int{}
	}
}

func (g *Graph) createEdges() {
	r := 60 * 60.0

	for i := range g.Nodes {
		p1 := g.Nodes[i].Position
		fmt.Println("p1", p1)
		for j := 0; j < i; j++ {
			p2 := g.Nodes[j].Position
			if Distance(p1, p2) <= r {
				g.AddEdge(i, j)
			}
		}

	}
}
