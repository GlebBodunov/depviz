package main

import (
	"flag"
	"fmt"
	"os"

	"depviz/src/apk"
	"depviz/src/config"
	"depviz/src/graph"
	"depviz/src/testrepo"
)

const defaultConfigPath = "config.yaml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "путь к YAML-файлу конфигурации")
	flag.Parse()

	if err := run(*configPath); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

func run(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	source, err := openSource(cfg)
	if err != nil {
		return err
	}
	g, err := graph.Build(source, cfg.PackageName, cfg.MaxDepth)
	if err != nil {
		return err
	}
	g.Print()
	return nil
}

func openSource(cfg *config.Config) (graph.Source, error) {
	if *cfg.TestMode {
		return testrepo.Load(cfg.Repository)
	}
	return apk.FetchIndex(cfg.Repository)
}
