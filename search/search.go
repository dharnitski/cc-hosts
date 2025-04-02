package search

import (
	"context"
	"fmt"
	"sort"
	"time"

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
	Target  string         `json:"target"`
	Out     []string       `json:"out"`
	In      []string       `json:"in"`
	Timings map[string]int `json:"timing"`
}

func (s *Searcher) GetTargets(ctx context.Context, domain string) (*Result, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain is empty")
	}
	reversed := vertices.ReverseDomain(domain)
	timings := make(map[string]int)
	start := time.Now()
	vertice, err := s.v.GetByDomain(ctx, reversed)
	if err != nil {
		return nil, err
	}
	timings["get_by_domain"] = int(time.Since(start).Milliseconds())
	if vertice == nil {
		return nil, nil
	}
	start = time.Now()
	outs, err := s.getDomains(ctx, vertice.ID(), s.out, timings, "out")
	if err != nil {
		return nil, err
	}
	timings["out_domains"] = int(time.Since(start).Milliseconds())
	start = time.Now()
	ins, err := s.getDomains(ctx, vertice.ID(), s.in, timings, "in")
	if err != nil {
		return nil, err
	}
	timings["in_domains"] = int(time.Since(start).Milliseconds())
	return &Result{Target: domain, Out: outs, In: ins, Timings: timings}, nil
}

func (s *Searcher) getDomains(ctx context.Context, verticeID string, edges *edges.Edges, timings map[string]int, pref string) ([]string, error) {
	start := time.Now()
	outIDs, err := edges.Get(ctx, verticeID)
	if err != nil {
		return nil, err
	}
	timings[fmt.Sprintf("edges_get_%s", pref)] = int(time.Since(start).Milliseconds())
	start = time.Now()
	domains, err := s.v.GetByIDs(ctx, outIDs)
	if err != nil {
		return nil, err
	}
	timings[fmt.Sprintf("v_get_by_ids_%s", pref)] = int(time.Since(start).Milliseconds())
	results := make([]string, 0, len(domains))
	for _, d := range domains {
		results = append(results, vertices.ReverseDomain(d.Domain()))
	}
	sort.Strings(results)
	return results, nil
}
