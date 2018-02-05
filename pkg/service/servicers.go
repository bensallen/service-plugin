package service

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/u-root/service-plugin/pkg/graph"
)

// Servicers type is a map which keys are the
// servicer name and their values are the actual Servicer
type Servicers struct {
	// Map of Servicers with names as keys
	Lookup map[string]Servicer

	// Sorted contains the resources after a topological sort.
	Sorted []*graph.Node

	// Reversed contains the resource dependency graph in reverse
	// order. It is used for finding the reverse dependencies of
	// services.
	Reversed *graph.Graph
}

// DependencyGraph builds a dependency graph for the collection
func (s Servicers) DependencyGraph() (*graph.Graph, error) {
	g := graph.New()

	// A map containing the resource ids and their nodes in the graph
	nodes := make(map[string]*graph.Node)
	for id := range s.Lookup {
		node := graph.NewNode(id)
		nodes[id] = node
		g.AddNode(node)
	}

	// Connect the nodes in the graph
	for id, svcr := range s.Lookup {

		u := svcr.Unit()

		// Create edges between the service and the ones
		// it wants to start after
		for _, dep := range u.After {
			if _, ok := s.Lookup[dep]; !ok {
				return g, fmt.Errorf("%s wants to start after %s, which does not exist", id, dep)
			}
			g.AddEdge(nodes[id], nodes[dep])

		}

		// Create edges between the service and the ones
		// it wants to start before
		for _, dep := range u.Before {
			if _, ok := s.Lookup[dep]; !ok {
				return g, fmt.Errorf("%s wants to start before %s, which does not exist", id, dep)
			}
			g.AddEdge(nodes[dep], nodes[id])
		}
	}

	return g, nil
}

func (s *Servicers) StartAll() error {

	// Start goroutines for concurrent processing
	var wg sync.WaitGroup
	ch := make(chan Servicer, 1024)

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		worker := func() {
			defer wg.Done()
			for svc := range ch {
				svc.Start()
			}
		}
		go worker()
	}

	// Process the Servicers
	for _, node := range s.Sorted {
		svcr := s.Lookup[node.Name]
		u := svcr.Unit()
		switch {
		// Servicer is an isolated node
		case len(u.After) == 0 && len(s.Reversed.Nodes[u.Name].Edges) == 0:
			ch <- svcr
			//fmt.Printf("Servicer: %s is an isolated node\n", u.Name)
			continue
		// Servicer is has no reverse dependencies
		case len(s.Reversed.Nodes[u.Name].Edges) == 0:
			ch <- svcr
			//fmt.Printf("Servicer: %s has no reverse dependencies\n", u.Name)
			continue
		// Servicer is not concurrent
		default:
			//fmt.Printf("Servicer: %s is not concurrent\n", u.Name)
			svcr.Start()
		}
	}

	close(ch)
	wg.Wait()

	return nil
}
