package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"depviz/src/apk"
	"depviz/src/config"
	"depviz/src/graph"
	"depviz/src/testrepo"
)

const defaultConfigPath = "config.yaml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "путь к YAML-файлу конфигурации")
	reverse := flag.Bool("reverse", false, "вывести обратные зависимости пакета")
	flag.Parse()

	if err := run(*configPath, *reverse); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

func run(configPath string, reverse bool) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	source, err := openSource(cfg)
	if err != nil {
		return err
	}
	if reverse {
		return printReverse(source, cfg)
	}
	g, err := graph.Build(source, cfg.PackageName, cfg.MaxDepth)
	if err != nil {
		return err
	}
	g.Print("Граф зависимостей", "->")
	return nil
}

func openSource(cfg *config.Config) (graph.Lister, error) {
	if *cfg.TestMode {
		return testrepo.Load(cfg.Repository)
	}
	return apk.FetchIndex(cfg.Repository)
}

func printReverse(source graph.Lister, cfg *config.Config) error {
	g, err := graph.Build(graph.Reverse(source), cfg.PackageName, cfg.MaxDepth)
	if err != nil {
		return err
	}
	g.Print("Обратные зависимости", "<-")
	reached := g.Reached()
	if len(reached) == 0 {
		fmt.Printf("От пакета %s никто не зависит\n", cfg.PackageName)
		return nil
	}
	fmt.Printf("От пакета %s зависят: %s\n", cfg.PackageName, strings.Join(reached, ", "))
	return nil
}
