package search

import (
	"context"
	"fmt"
	"sort"

	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/vertices"
)

type Searcher struct {
	// from target to other sites
	out *edges.Edges
	// from external sites to target
	in *edges.Edges
	v  *vertices.Vertices
}

func NewSearcher(v *vertices.Vertices, out *edges.Edges, in *edges.Edges) *Searcher {
	return &Searcher{v: v, out: out, in: in}
}

type Result struct {
	Target string   `json:"target"`
	Out    []string `json:"out"`
	In     []string `json:"in"`
}

func (s *Searcher) GetTargets(ctx context.Context, domain string) (*Result, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain is empty")
	}
	reversed := vertices.ReverseDomain(domain)
	vertice, err := s.v.GetByDomain(ctx, reversed)
	if err != nil {
		return nil, err
	}
	if vertice == nil {
		return nil, nil
	}
	outs, err := s.getDomains(ctx, vertice.ID(), s.out)
	if err != nil {
		return nil, err
	}
	ins, err := s.getDomains(ctx, vertice.ID(), s.in)
	if err != nil {
		return nil, err
	}
	return &Result{Target: domain, Out: outs, In: ins}, nil
}

func (s *Searcher) getDomains(ctx context.Context, verticeID string, edges *edges.Edges) ([]string, error) {
	outIDs, err := edges.Get(ctx, verticeID)
	if err != nil {
		return nil, err
	}
	domains, err := s.v.GetByIDs(ctx, outIDs)
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
