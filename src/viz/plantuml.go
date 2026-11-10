// Package viz формирует описание графа на языке PlantUML и сохраняет изображение графа в PNG.
package viz

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"depviz/src/graph"
)

const (
	plantumlCmd = "plantuml"
	rootColor   = "#LightBlue"
	edgeArrow   = "-->"
	cycleArrow  = "-[#red]->"
	sourceExt   = ".puml"
	filePerm    = 0o644
)

func PlantUML(g *graph.Graph) string {
	cycles := map[graph.Edge]bool{}
	for _, edge := range g.Cycles {
		cycles[edge] = true
	}
	var sb strings.Builder
	sb.WriteString("@startuml\n")
	fmt.Fprintf(&sb, "[%s] %s\n", g.Root, rootColor)
	for _, node := range g.Nodes {
		for _, dep := range g.Edges[node] {
			arrow := edgeArrow
			if cycles[graph.Edge{From: node, To: dep}] {
				arrow = cycleArrow
			}
			fmt.Fprintf(&sb, "[%s] %s [%s]\n", node, arrow, dep)
		}
	}
	sb.WriteString("@enduml\n")
	return sb.String()
}

func SourcePath(imagePath string) string {
	return strings.TrimSuffix(imagePath, filepath.Ext(imagePath)) + sourceExt
}

func Save(source, imagePath string) error {
	if err := os.WriteFile(SourcePath(imagePath), []byte(source), filePerm); err != nil {
		return fmt.Errorf("не удалось сохранить описание графа: %w", err)
	}
	if _, err := exec.LookPath(plantumlCmd); err != nil {
		return fmt.Errorf("программа %s не найдена, установите её: brew install plantuml", plantumlCmd)
	}
	cmd := exec.Command(plantumlCmd, "-tpng", "-pipe")
	cmd.Stdin = strings.NewReader(source)
	var image, stderr bytes.Buffer
	cmd.Stdout = &image
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("PlantUML завершился с ошибкой: %w: %s", err, stderr.String())
	}
	if err := os.WriteFile(imagePath, image.Bytes(), filePerm); err != nil {
		return fmt.Errorf("не удалось сохранить изображение: %w", err)
	}
	return nil
}
