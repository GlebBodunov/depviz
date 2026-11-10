package tests

import (
	"reflect"
	"testing"

	"depviz/src/graph"
	"depviz/src/testrepo"
)

func buildReverse(t *testing.T, file, root string, depth int) *graph.Graph {
	t.Helper()
	repo, err := testrepo.Load("testdata/" + file)
	if err != nil {
		t.Fatalf("не удалось загрузить %s: %v", file, err)
	}
	g, err := graph.Build(graph.Reverse(repo), root, depth)
	if err != nil {
		t.Fatalf("не удалось построить обратный граф: %v", err)
	}
	return g
}

func TestReverseDiamond(t *testing.T) {
	g := buildReverse(t, "simple.txt", "D", unlimitedDepth)
	want := []string{"B", "A", "C"}
	if !reflect.DeepEqual(g.Reached(), want) {
		t.Errorf("обратные зависимости D: %v, ожидалось %v", g.Reached(), want)
	}
}

func TestReverseRootHasNoDependents(t *testing.T) {
	g := buildReverse(t, "simple.txt", "A", unlimitedDepth)
	if len(g.Reached()) != 0 {
		t.Errorf("от A никто не зависит, получено %v", g.Reached())
	}
}

func TestReverseCycle(t *testing.T) {
	g := buildReverse(t, "cycle.txt", "A", unlimitedDepth)
	if !reflect.DeepEqual(g.Reached(), []string{"C", "B"}) {
		t.Errorf("обратные зависимости A: %v, ожидалось [C B]", g.Reached())
	}
	if len(g.Cycles) != 1 {
		t.Errorf("ожидался один цикл, найдено %v", g.Cycles)
	}
}

func TestReverseMaxDepth(t *testing.T) {
	g := buildReverse(t, "deep.txt", "E", 2)
	if !reflect.DeepEqual(g.Reached(), []string{"D", "C"}) {
		t.Errorf("обратные зависимости E на глубине 2: %v, ожидалось [D C]", g.Reached())
	}
}

func TestReverseUnknownPackage(t *testing.T) {
	repo, _ := testrepo.Load("testdata/simple.txt")
	if _, err := graph.Build(graph.Reverse(repo), "Z", unlimitedDepth); err == nil {
		t.Error("ожидалась ошибка для несуществующего пакета")
	}
}
