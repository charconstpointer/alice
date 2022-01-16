package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/charconstpointer/alice/pr"
)

func main() {
	app := App{}
	src := pr.NewSource(http.Client{Timeout: 5 * time.Second})
	err := app.AddSource(src)
	if err != nil {
		log.Fatalf("failed to add source: %v", err)
	}

	for {
		// scan os stdin input
		var input string
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()
		results, err := app.Find(input)
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

func (a *App) AddSource(source Source) error {
	for _, s := range a.sources {
		if s == source {
			return fmt.Errorf("source already exists")
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

	return results, nil
}
