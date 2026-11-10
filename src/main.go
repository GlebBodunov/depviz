package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"depviz/src/apk"
	"depviz/src/config"
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
	if *cfg.TestMode {
		return errors.New("тестовый режим будет реализован на этапе 3")
	}
	idx, err := apk.FetchIndex(cfg.Repository)
	if err != nil {
		return err
	}
	deps, err := idx.DirectDeps(cfg.PackageName)
	if err != nil {
		return err
	}
	fmt.Printf("Прямые зависимости пакета %s:\n", cfg.PackageName)
	for _, dep := range deps {
		fmt.Println("  " + dep)
	}
	return nil
}
