// Package testrepo читает тестовый репозиторий, где пакеты названы большими латинскими буквами.
package testrepo

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const commentPrefix = "#"

var namePattern = regexp.MustCompile(`^[A-Z]+$`)

type Repo struct {
	deps map[string][]string
}

func Load(path string) (*Repo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть тестовый репозиторий: %w", err)
	}
	defer file.Close()
	repo := &Repo{deps: map[string][]string{}}
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, commentPrefix) {
			continue
		}
		if err := repo.parseLine(line); err != nil {
			return nil, fmt.Errorf("строка %d: %w", lineNum, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения тестового репозитория: %w", err)
	}
	return repo, nil
}

func (r *Repo) parseLine(line string) error {
	name, rest, ok := strings.Cut(line, ":")
	if !ok {
		return fmt.Errorf("ожидается формат «ПАКЕТ: ЗАВИСИМОСТИ», получено %q", line)
	}
	name = strings.TrimSpace(name)
	deps := strings.Fields(rest)
	for _, n := range append([]string{name}, deps...) {
		if !namePattern.MatchString(n) {
			return fmt.Errorf("имя пакета %q должно состоять из больших латинских букв", n)
		}
	}
	if _, exists := r.deps[name]; exists {
		return fmt.Errorf("пакет %s описан повторно", name)
	}
	r.deps[name] = deps
	return nil
}

func (r *Repo) DirectDeps(name string) ([]string, error) {
	deps, ok := r.deps[name]
	if !ok {
		return nil, fmt.Errorf("пакет %q не найден в тестовом репозитории", name)
	}
	return deps, nil
}
