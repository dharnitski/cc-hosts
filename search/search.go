package search

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	real_aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/access/aws"
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
	mu sync.Mutex
}

func NewSearcher(ctx context.Context, offsetsGetter access.Getter, edgesGetter access.Getter, revEdgesGetter access.Getter, verticesGetter access.Getter) (*Searcher, error) {
	eOffsets, err := edges.NewOffsets(ctx, offsetsGetter)
	if err != nil {
		return nil, err
	}

	out := edges.NewEdges(edgesGetter, *eOffsets)

	reversedOffsets, err := edges.NewOffsetsReversed(ctx, offsetsGetter)
	if err != nil {
		return nil, err
	}

	in := edges.NewEdges(revEdgesGetter, *reversedOffsets)

	vOffsets, err := vertices.NewOffsets(ctx, offsetsGetter)
	if err != nil {
		return nil, err
	}

	v := vertices.NewVertices(verticesGetter, *vOffsets)

	searcher := &Searcher{v: v, out: out, in: in}

	return searcher, nil
}

func NewAwsSearcher(ctx context.Context, cfg real_aws.Config, bucket string) (*Searcher, error) {
	// folder is empty - expect offset files in the root of the bucket
	offsetsGetter := aws.New(cfg, bucket, "")
	edgesGetter := aws.New(cfg, bucket, edges.EdgesFolder)
	revEdgesGetter := aws.New(cfg, bucket, edges.EdgesReversedFolder)
	verticesGetter := aws.New(cfg, bucket, vertices.Folder)

	return NewSearcher(ctx, offsetsGetter, edgesGetter, revEdgesGetter, verticesGetter)
}

type Result struct {
	Target  string         `json:"target"`
	Out     []string       `json:"out"`
	In      []string       `json:"in"`
	Timings map[string]int `json:"timing"`
}

func (s *Searcher) GetTargets(ctx context.Context, domain string) (*Result, error) {
	if domain == "" {
		return nil, errors.New("domain is empty")
	}

	reversed := vertices.ReverseDomain(domain)
	timings := make(map[string]int)
	start := time.Now()

	vertex, err := s.v.GetByDomain(ctx, reversed)
	if err != nil {
		return nil, err
	}

	timings["get_by_domain"] = int(time.Since(start).Milliseconds())

	if vertex == nil {
		return nil, nil //nolint:nilnil
	}

	var outs, ins []string

	var outErr, inErr error

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		outs, outErr = s.getDomains(ctx, vertex.ID(), timings, out)
	}()

	go func() {
		defer wg.Done()

		ins, inErr = s.getDomains(ctx, vertex.ID(), timings, in)
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

	s.mu.Lock()
	timings[fmt.Sprintf("edges_get_%s", pref)] = int(time.Since(start).Milliseconds())
	s.mu.Unlock()

	start = time.Now()

	domains, err := s.v.GetByIDs(ctx, outIDs)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	timings[fmt.Sprintf("v_get_by_ids_%s", pref)] = int(time.Since(start).Milliseconds())
	s.mu.Unlock()

	results := make([]string, 0, len(domains))
	for _, d := range domains {
		results = append(results, vertices.ReverseDomain(d.Domain()))
	}

	sort.Strings(results)

	s.mu.Lock()
	timings[fmt.Sprintf("%s_domains", pref)] = int(time.Since(allStart).Milliseconds())
	s.mu.Unlock()

	return results, nil
}
