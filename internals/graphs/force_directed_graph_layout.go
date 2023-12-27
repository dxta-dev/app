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
	k := 350.0
	return k * k / distance
}

func attractiveForce(distance float64) float64 {
	if distance == 0 {
		distance = 0.1
	}
	k := 100.0
	return distance * distance / k
}

func printGraph(iteration int, graph *Graph) {
	for i := 0; i < len(graph.Nodes); i++ {
		fmt.Println(iteration, graph.Nodes[i].Position.X, graph.Nodes[i].Position.Y)
	}
}

func ForceDirectedGraphLayout(graph *Graph, iterations int) {
	temperature := 60 * 15.0

	for iter := 0; iter < iterations; iter++ {
		printGraph(iter, graph)
		displacement := make([]Point, len(graph.Nodes))

		/*for i, v := range graph.Nodes {
			for j, u := range graph.Nodes {
				if i != j {
					delta := Point{X: v.Position.X - u.Position.X, Y: v.Position.Y - u.Position.Y}
					distance := math.Sqrt(delta.X*delta.X+delta.Y*delta.Y)
					force := repulsiveForce(distance)
					fmt.Println("repu force", i, j, force, distance)
					fmt.Println(delta.Y, distance)
					displacement[i].Y += force;
				}
			}
		}*/

		for _, edge := range graph.Edges {
			v := graph.Nodes[edge[0]]
			u := graph.Nodes[edge[1]]
			delta := Point{X: v.Position.X - u.Position.X, Y: v.Position.Y - u.Position.Y}
			distance := math.Sqrt(delta.X*delta.X + delta.Y*delta.Y)
			force := repulsiveForce(distance)

			if delta.Y > 0 {
				displacement[edge[0]].Y += force
				displacement[edge[1]].Y -= force
			} else {
				displacement[edge[0]].Y -= force
				displacement[edge[1]].Y += force
			}
		}

		/*for _, edge := range graph.Edges {
			v := graph.Nodes[edge[0]]
			u := graph.Nodes[edge[1]]
			delta := Point{X: v.Position.X - u.Position.X, Y: v.Position.Y - u.Position.Y}
			distance := math.Sqrt(delta.X*delta.X + delta.Y*delta.Y)
			force := attractiveForce(distance)
			fmt.Println("attr force", edge[0], edge[1], force)
			displacement[edge[0]].Y -= force;
			displacement[edge[1]].Y += force;
		}*/

		for i, v := range graph.Nodes {
			dispLength := math.Abs(displacement[i].Y)
			fmt.Println(displacement[i].Y / dispLength * math.Min(dispLength, temperature))
			if displacement[i].Y != 0 {
				v.Position.Y += displacement[i].Y / dispLength * math.Min(dispLength, temperature)
			}
		}

		temperature *= 0.95
	}
}

func FindClosePoints(index int, xValues []float64, yValues []float64, p1 Point, r float64) []int {
	var closePoints []int
	for i := range xValues {
		if i == index {
			continue
		}
		p := Point{X: xValues[i], Y: yValues[i]}
		if Distance(p1, p) <= r {
			closePoints = append(closePoints, i)
		}
	}
	return closePoints
}
