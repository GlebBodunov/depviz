// Package tests содержит тесты для пакетов приложения depviz.
package tests

import (
	"strings"
	"testing"

	"depviz/src/config"
)

const validYAML = `
package_name: busybox
repository: https://dl-cdn.alpinelinux.org/alpine/v3.20/main/x86_64/APKINDEX.tar.gz
test_mode: false
output_file: graph.png
max_depth: 3
`

func TestParseValid(t *testing.T) {
	cfg, err := config.Parse([]byte(validYAML))
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if cfg.PackageName != "busybox" {
		t.Errorf("package_name = %q, ожидалось busybox", cfg.PackageName)
	}
}

func TestParseErrors(t *testing.T) {
	cases := map[string]string{
		"package_name": strings.Replace(validYAML, "busybox", `""`, 1),
		"test_mode":    strings.Replace(validYAML, "test_mode: false\n", "", 1),
		"repository":   strings.Replace(validYAML, "https://", "ftp://", 1),
		"output_file":  strings.Replace(validYAML, "graph.png", "graph.jpg", 1),
		"max_depth":    strings.Replace(validYAML, "max_depth: 3", "max_depth: 0", 1),
		"unknown":      validYAML + "color: red\n",
		"syntax":       "package_name: [busybox",
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := config.Parse([]byte(data)); err == nil {
				t.Errorf("ожидалась ошибка для случая %s", name)
			}
		})
	}
}

func TestTestModeMissingFile(t *testing.T) {
	data := strings.Replace(validYAML, "test_mode: false", "test_mode: true", 1)
	if _, err := config.Parse([]byte(data)); err == nil {
		t.Error("ожидалась ошибка: тестовый файл не существует")
	}
}
