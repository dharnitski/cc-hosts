package search

import (
	"context"
	"sort"

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

func (s *Searcher) GetTargets(ctx context.Context, domain string) ([]string, error) {
	domain = vertices.ReverseDomain(domain)
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
		results = append(results, vertices.ReverseDomain(d.Domain()))
	}
	sort.Strings(results)
	return results, nil
}
