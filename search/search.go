package search

import (
	"context"
	"sort"
	"strings"

	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/vertices"
)

type Searcher struct {
	e *edges.Edges
	v *vertices.Vertices
}

func NewSearcher(v *vertices.Vertices, e *edges.Edges) *Searcher {
	return &Searcher{v: v, e: e}
}

func reverseDomain(domain string) string {
	parts := strings.Split(domain, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, ".")
}

func (s *Searcher) GetTargets(ctx context.Context, domain string) ([]string, error) {
	domain = reverseDomain(domain)
	vertice, err := s.v.GetByDomain(domain)
	if err != nil {
		return nil, err
	}
	if vertice == nil {
		return nil, nil
	}
	eIDs, err := s.e.Get(vertice.ID())
	if err != nil {
		return nil, err
	}
	domains, err := s.v.GetByIDs(eIDs)
	if err != nil {
		return nil, err
	}
	results := make([]string, 0, len(domains))
	for _, d := range domains {
		results = append(results, reverseDomain(d.Domain()))
	}
	sort.Strings(results)
	return results, nil
}
