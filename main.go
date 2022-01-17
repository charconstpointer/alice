package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/charconstpointer/alice/pr"
	"github.com/charconstpointer/alice/wiki"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func main() {
	app := App{}
	prSrc := pr.NewSource(http.Client{Timeout: 5 * time.Second})
	wikiSrc := wiki.NewSource(http.Client{Timeout: 5 * time.Second})

	err := app.AddSources(prSrc, wikiSrc)
	if err != nil {
		log.Fatalf("failed to add source: %v", err)
	}

	for {
		// scan os stdin input
		// var input string
		// scanner := bufio.NewScanner(os.Stdin)
		// scanner.Scan()
		// input = scanner.Text()
		results, err := app.Find("kosa")
		if err != nil {
			log.Printf("failed to find: %v", err)
			continue
		}
		for _, result := range results {
			fmt.Printf(">>> %s\n", result)
		}
	}
}

type App struct {
	sources []Source
}

type Source interface {
	Find(context.Context, string) ([]string, error)
}

func (a *App) AddSources(source ...Source) error {
	for _, s := range source {
		err := a.AddSource(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) AddSource(source Source) error {
	//check if source is already added
	for _, s := range a.sources {
		if s == source {
			return fmt.Errorf("source already added")
		}
	}

	a.sources = append(a.sources, source)
	return nil
}

func (a *App) Find(phrase string) ([]string, error) {
	resCh := make(chan string)
	var results []string
	var wg sync.WaitGroup
	for _, s := range a.sources {
		wg.Add(1)
		go func(s Source) {
			defer wg.Done()
			results, err := s.Find(context.Background(), phrase)
			if err != nil {
				log.Printf("failed to find: %v", err)
			}
			for _, r := range results {
				resCh <- r
			}
		}(s)
	}
	go func() {
		wg.Wait()
		close(resCh)
	}()
	for r := range resCh {
		results = append(results, r)
	}

	final := fuzzy.Find(phrase, results)
	if len(final) == 0 {
		return results, fmt.Errorf("no results")
	}
	return final[:3], nil
}
