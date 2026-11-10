// Package graph строит граф зависимостей рекурсивным обходом в глубину (DFS) с ограничением глубины.
package graph

import "fmt"

type Source interface {
	DirectDeps(name string) ([]string, error)
}

type Edge struct {
	From string
	To   string
}

type Graph struct {
	Root    string
	Nodes   []string
	Edges   map[string][]string
	Cycles  []Edge
	Missing []string
}

type builder struct {
	source   Source
	maxDepth int
	graph    *Graph
	onStack  map[string]bool
	depthOf  map[string]int
	edgeSeen map[Edge]bool
}

func Build(source Source, root string, maxDepth int) (*Graph, error) {
	if _, err := source.DirectDeps(root); err != nil {
		return nil, err
	}
	b := &builder{
		source:   source,
		maxDepth: maxDepth,
		graph:    &Graph{Root: root, Edges: map[string][]string{}},
		onStack:  map[string]bool{},
		depthOf:  map[string]int{},
		edgeSeen: map[Edge]bool{},
	}
	b.visit(root, 0)
	return b.graph, nil
}

func (b *builder) visit(name string, depth int) {
	if prev, seen := b.depthOf[name]; seen && prev <= depth {
		return
	}
	if _, seen := b.depthOf[name]; !seen {
		b.graph.Nodes = append(b.graph.Nodes, name)
	}
	b.depthOf[name] = depth
	if depth >= b.maxDepth {
		return
	}
	deps, err := b.source.DirectDeps(name)
	if err != nil {
		b.graph.Missing = append(b.graph.Missing, name)
		return
	}
	b.onStack[name] = true
	for _, dep := range deps {
		b.addEdge(name, dep)
		if b.onStack[dep] {
			b.graph.Cycles = append(b.graph.Cycles, Edge{From: name, To: dep})
			continue
		}
		b.visit(dep, depth+1)
	}
	b.onStack[name] = false
}

func (b *builder) addEdge(from, to string) {
	edge := Edge{From: from, To: to}
	if b.edgeSeen[edge] {
		return
	}
	b.edgeSeen[edge] = true
	b.graph.Edges[from] = append(b.graph.Edges[from], to)
}

func (g *Graph) Print() {
	fmt.Printf("Граф зависимостей пакета %s:\n", g.Root)
	for _, node := range g.Nodes {
		if deps := g.Edges[node]; len(deps) > 0 {
			fmt.Printf("  %s -> %v\n", node, deps)
		}
	}
	for _, cycle := range g.Cycles {
		fmt.Printf("Обнаружен цикл: %s -> %s\n", cycle.From, cycle.To)
	}
	for _, name := range g.Missing {
		fmt.Printf("Пакет %s не найден в репозитории\n", name)
	}
}
