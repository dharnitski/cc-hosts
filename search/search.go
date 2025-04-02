package search

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/vertices"
)

type direction string

const (
	in  direction = "in"
	out direction = "out"
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
	var outs, ins []string
	var outErr, inErr error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		outs, outErr = s.getDomains(ctx, vertice.ID(), timings, out)
	}()

	go func() {
		defer wg.Done()
		ins, inErr = s.getDomains(ctx, vertice.ID(), timings, in)
	}()
	wg.Wait()

	if outErr != nil {
		return nil, outErr
	}
	if inErr != nil {
		return nil, inErr
	}
	return &Result{Target: domain, Out: outs, In: ins, Timings: timings}, nil
}

func (s *Searcher) getDomains(ctx context.Context, verticeID string, timings map[string]int, pref direction) ([]string, error) {
	allStart := time.Now()
	var edges *edges.Edges
	switch pref {
	case out:
		edges = s.out
	case in:
		edges = s.in
	}

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
	timings[fmt.Sprintf("%s_domains", pref)] = int(time.Since(allStart).Milliseconds())
	return results, nil
}
