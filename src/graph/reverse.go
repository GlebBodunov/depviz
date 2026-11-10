package graph

import "fmt"

type Lister interface {
	Source
	Names() []string
}

type reverseSource struct {
	dependents map[string][]string
	known      map[string]bool
}

func Reverse(source Lister) Source {
	rev := &reverseSource{dependents: map[string][]string{}, known: map[string]bool{}}
	for _, name := range source.Names() {
		rev.known[name] = true
		deps, err := source.DirectDeps(name)
		if err != nil {
			continue
		}
		for _, dep := range deps {
			rev.known[dep] = true
			rev.dependents[dep] = append(rev.dependents[dep], name)
		}
	}
	return rev
}

func (r *reverseSource) DirectDeps(name string) ([]string, error) {
	if !r.known[name] {
		return nil, fmt.Errorf("пакет %q не найден в репозитории", name)
	}
	return r.dependents[name], nil
}
