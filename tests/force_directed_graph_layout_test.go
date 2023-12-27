package main

import (
	"dxta-dev/app/internals/graphs"
	"fmt"
	"testing"
)

// Add adds two integers
func Add(a, b int) int {
	return a + b
}

// TestAdd tests the Add function
func TestAdd(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", result)
	}
}
func TestSimpleGraph(t *testing.T) {

	points := [][2]float64{
		{60*60*12, 60*60*24},
		{60*60*12, 60*60*24},
	}

	graph := graphs.NewGraph()
	for i := 0; i < len(points); i++ {
		graph.AddNode(graphs.NewNode(points[i][0], points[i][1]))
		/*for i := 0; i < len(graph.Nodes); i++ {
			closePoints := FindClosePoints(i, xValues, yValues, graph.Nodes[i].Position, 60*60*6)
			for _, p := range closePoints {
				newGraph.AddEdge(i, p)
			}
		}*/
	}

	graph.AddEdge(0, 1)

	graphs.ForceDirectedGraphLayout(graph, 1000)

	for i := 0; i < len(graph.Nodes); i++ {
		fmt.Println("Result:", i, graph.Nodes[i].Position.X, graph.Nodes[i].Position.Y)
	}

	t.Errorf("SimpleGraph")
}
