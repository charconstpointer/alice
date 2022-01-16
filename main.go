package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	var results []string
	for _, s := range a.sources {
		res, err := s.Find(context.Background(), phrase)
		if err != nil {
			return nil, err
		}
		results = append(results, res...)
	}
	return results, nil
}
