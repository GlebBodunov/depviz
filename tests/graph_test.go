package tests

import (
	"reflect"
	"testing"

	"depviz/src/graph"
	"depviz/src/testrepo"
)

const unlimitedDepth = 100

func buildGraph(t *testing.T, file, root string, depth int) *graph.Graph {
	t.Helper()
	repo, err := testrepo.Load("testdata/" + file)
	if err != nil {
		t.Fatalf("не удалось загрузить %s: %v", file, err)
	}
	g, err := graph.Build(repo, root, depth)
	if err != nil {
		t.Fatalf("не удалось построить граф: %v", err)
	}
	return g
}

func TestSimpleGraph(t *testing.T) {
	g := buildGraph(t, "simple.txt", "A", unlimitedDepth)
	want := []string{"A", "B", "D", "C"}
	if !reflect.DeepEqual(g.Nodes, want) {
		t.Errorf("порядок обхода %v, ожидалось %v", g.Nodes, want)
	}
	if len(g.Cycles) != 0 {
		t.Errorf("циклов быть не должно, найдено %v", g.Cycles)
	}
}

func TestCycle(t *testing.T) {
	g := buildGraph(t, "cycle.txt", "A", unlimitedDepth)
	want := []graph.Edge{{From: "C", To: "A"}}
	if !reflect.DeepEqual(g.Cycles, want) {
		t.Errorf("циклы %v, ожидалось %v", g.Cycles, want)
	}
}

func TestMaxDepth(t *testing.T) {
	cases := map[int][]string{
		1: {"A", "B"},
		2: {"A", "B", "C"},
		4: {"A", "B", "C", "D", "E"},
	}
	for depth, want := range cases {
		g := buildGraph(t, "deep.txt", "A", depth)
		if !reflect.DeepEqual(g.Nodes, want) {
			t.Errorf("глубина %d: узлы %v, ожидалось %v", depth, g.Nodes, want)
		}
	}
}

func TestShallowerPathExpandsNode(t *testing.T) {
	g := buildGraph(t, "shortcut.txt", "A", 2)
	want := []string{"A", "B", "C", "D"}
	if !reflect.DeepEqual(g.Nodes, want) {
		t.Errorf("узлы %v, ожидалось %v", g.Nodes, want)
	}
}

func TestComplexGraph(t *testing.T) {
	g := buildGraph(t, "complex.txt", "A", unlimitedDepth)
	wantCycles := []graph.Edge{{From: "D", To: "B"}, {From: "G", To: "A"}, {From: "F", To: "F"}}
	if !reflect.DeepEqual(g.Cycles, wantCycles) {
		t.Errorf("циклы %v, ожидалось %v", g.Cycles, wantCycles)
	}
	if !reflect.DeepEqual(g.Missing, []string{"X"}) {
		t.Errorf("отсутствующие пакеты %v, ожидалось [X]", g.Missing)
	}
}

func TestUnknownRoot(t *testing.T) {
	repo, _ := testrepo.Load("testdata/simple.txt")
	if _, err := graph.Build(repo, "Z", unlimitedDepth); err == nil {
		t.Error("ожидалась ошибка для несуществующего корневого пакета")
	}
}

func TestInvalidTestRepo(t *testing.T) {
	if _, err := testrepo.Load("testdata/invalid.txt"); err == nil {
		t.Error("ожидалась ошибка для имени из маленьких букв")
	}
}
