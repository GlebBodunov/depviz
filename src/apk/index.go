// Package apk извлекает информацию о пакетах из индекса репозитория Alpine Linux.
package apk

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const versionOperators = "<>=~"

const conflictPrefix = "!"

type Package struct {
	Name     string
	Version  string
	Depends  []string
	Provides []string
}

type Index struct {
	Packages  map[string]*Package
	providers map[string]string
}

func ParseIndex(r io.Reader) (*Index, error) {
	idx := &Index{Packages: map[string]*Package{}, providers: map[string]string{}}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, bufio.MaxScanTokenSize), bufio.MaxScanTokenSize*16)
	current := &Package{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			idx.add(current)
			current = &Package{}
			continue
		}
		parseLine(current, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения APKINDEX: %w", err)
	}
	idx.add(current)
	return idx, nil
}

func parseLine(pkg *Package, line string) {
	key, value, ok := strings.Cut(line, ":")
	if !ok {
		return
	}
	switch key {
	case "P":
		pkg.Name = value
	case "V":
		pkg.Version = value
	case "D":
		pkg.Depends = strings.Fields(value)
	case "p":
		pkg.Provides = strings.Fields(value)
	}
}

func (idx *Index) add(pkg *Package) {
	if pkg.Name == "" {
		return
	}
	idx.Packages[pkg.Name] = pkg
	for _, provided := range pkg.Provides {
		idx.providers[stripVersion(provided)] = pkg.Name
	}
}

func stripVersion(dep string) string {
	if pos := strings.IndexAny(dep, versionOperators); pos >= 0 {
		return dep[:pos]
	}
	return dep
}

func (idx *Index) resolve(dep string) string {
	name := stripVersion(dep)
	if _, ok := idx.Packages[name]; ok {
		return name
	}
	if provider, ok := idx.providers[name]; ok {
		return provider
	}
	return name
}

func (idx *Index) DirectDeps(name string) ([]string, error) {
	pkg, ok := idx.Packages[name]
	if !ok {
		return nil, fmt.Errorf("пакет %q не найден в репозитории", name)
	}
	seen := map[string]bool{}
	var deps []string
	for _, dep := range pkg.Depends {
		if strings.HasPrefix(dep, conflictPrefix) {
			continue
		}
		resolved := idx.resolve(dep)
		if resolved == name || seen[resolved] {
			continue
		}
		seen[resolved] = true
		deps = append(deps, resolved)
	}
	return deps, nil
}
