// Package config загружает и проверяет параметры приложения из YAML-файла.
package config

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const minDepth = 1
const imageExt = ".png"

type Config struct {
	PackageName string `yaml:"package_name"`
	Repository  string `yaml:"repository"`
	TestMode    *bool  `yaml:"test_mode"`
	OutputFile  string `yaml:"output_file"`
	MaxDepth    int    `yaml:"max_depth"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
	}
	return Parse(data)
}

func Parse(data []byte) (*Config, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка синтаксиса YAML: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	checks := []func() error{
		c.validatePackageName,
		c.validateTestMode,
		c.validateRepository,
		c.validateOutputFile,
		c.validateMaxDepth,
	}
	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) validatePackageName() error {
	if strings.TrimSpace(c.PackageName) == "" {
		return errors.New("package_name: имя пакета не задано")
	}
	if strings.ContainsAny(c.PackageName, " \t") {
		return errors.New("package_name: имя пакета не должно содержать пробелов")
	}
	return nil
}

func (c *Config) validateTestMode() error {
	if c.TestMode == nil {
		return errors.New("test_mode: параметр не задан (ожидается true или false)")
	}
	return nil
}

func (c *Config) validateRepository() error {
	if strings.TrimSpace(c.Repository) == "" {
		return errors.New("repository: адрес репозитория не задан")
	}
	if *c.TestMode {
		return checkFileExists(c.Repository)
	}
	return checkURL(c.Repository)
}

func checkFileExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("repository: тестовый файл недоступен: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("repository: %q — это папка, а не файл", path)
	}
	return nil
}

func checkURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("repository: некорректный URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("repository: схема URL должна быть http или https, получено %q", u.Scheme)
	}
	if u.Host == "" {
		return errors.New("repository: в URL отсутствует имя хоста")
	}
	return nil
}

func (c *Config) validateOutputFile() error {
	if strings.TrimSpace(c.OutputFile) == "" {
		return errors.New("output_file: имя файла изображения не задано")
	}
	if !strings.HasSuffix(strings.ToLower(c.OutputFile), imageExt) {
		return fmt.Errorf("output_file: файл должен иметь расширение %s", imageExt)
	}
	return nil
}

func (c *Config) validateMaxDepth() error {
	if c.MaxDepth < minDepth {
		return fmt.Errorf("max_depth: глубина должна быть не меньше %d, получено %d", minDepth, c.MaxDepth)
	}
	return nil
}
