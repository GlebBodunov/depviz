package main

import (
	"flag"
	"fmt"
	"os"

	"depviz/src/config"
)

// defaultConfigPath — путь к файлу конфигурации по умолчанию.
const defaultConfigPath = "config.yaml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "путь к YAML-файлу конфигурации")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
	cfg.Print()
}
