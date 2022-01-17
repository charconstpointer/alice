package wiki

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Source struct {
	c http.Client
}

func NewSource(c http.Client) *Source {
	return &Source{c}
}
func (s *Source) Find(ctx context.Context, phrase string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://en.wikipedia.org/wiki/Poland", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	res, err := s.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %v", err)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	r := bytes.NewReader(b)
	sc := bufio.NewScanner(r)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	reuslts := fuzzy.Find(phrase, lines)
	if len(reuslts) == 0 {
		return nil, fmt.Errorf("no results")
	}

	return reuslts, nil
}
