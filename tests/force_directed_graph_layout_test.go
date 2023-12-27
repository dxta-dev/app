package main

import (
	"dxta-dev/app/internals/graphs"
	"fmt"
	"testing"
)

func TestSimpleGraph(t *testing.T) {

	points := [][2]float64{
		{60*60*12, 60*60*24},
		{60*60*12.25, 60*60*24},
		{60*60*12.50, 60*60*24},
	}

	graph := graphs.NewGraph()

	graph.AddNode(graphs.NewNode(points[0][0], points[0][1]))
	graph.AddNode(graphs.NewNode(points[1][0], points[1][1]))
	graph.AddNode(graphs.NewNode(points[2][0], points[2][1]))

	graph.AddEdge(0, 1)
	graph.AddEdge(0, 2)
	graph.AddEdge(1, 2)

	graphs.ForceDirectedGraphLayout(graph, 1)

	for i := 0; i < len(graph.Nodes); i++ {
		fmt.Println("Result:", i, graph.Nodes[i].Position.X, graph.Nodes[i].Position.Y)
	}

	t.Errorf("SimpleGraph")
}
