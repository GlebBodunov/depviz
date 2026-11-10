package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"depviz/src/viz"
)

func TestPlantUMLText(t *testing.T) {
	g := buildGraph(t, "cycle.txt", "A", unlimitedDepth)
	got := viz.PlantUML(g)
	want := "@startuml\n[A] #LightBlue\n[A] --> [B]\n[B] --> [C]\n[C] -[#red]-> [A]\n@enduml\n"
	if got != want {
		t.Errorf("получено:\n%s\nожидалось:\n%s", got, want)
	}
}

func TestPlantUMLAllEdges(t *testing.T) {
	g := buildGraph(t, "simple.txt", "A", unlimitedDepth)
	text := viz.PlantUML(g)
	for _, edge := range []string{"[A] --> [B]", "[A] --> [C]", "[B] --> [D]", "[C] --> [D]"} {
		if !strings.Contains(text, edge) {
			t.Errorf("в описании нет ребра %s", edge)
		}
	}
}

func TestSourcePath(t *testing.T) {
	if got := viz.SourcePath("out/graph.png"); got != "out/graph.puml" {
		t.Errorf("получено %s, ожидалось out/graph.puml", got)
	}
}

func TestSaveWritesSource(t *testing.T) {
	image := filepath.Join(t.TempDir(), "graph.png")
	_ = viz.Save("@startuml\n[A]\n@enduml\n", image)
	if _, err := os.Stat(viz.SourcePath(image)); err != nil {
		t.Errorf("файл с описанием графа не создан: %v", err)
	}
}
