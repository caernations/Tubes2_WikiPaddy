package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)


// PAGE LINKS (UNTUK NGESTORE & NGEMANAGE SEMUA DAFTAR LINK YANG ADA)
// PageLinks represents a mapping from a page to its links
type PageLinks struct {
	mu    sync.Mutex
	links map[string][]string
}

// NewPageLinks creates a new PageLinks instance
func NewPageLinks() *PageLinks {
	return &PageLinks{
		links: make(map[string][]string),
	}
}

// Add adds a link to the page
func (pl *PageLinks) Add(page, link string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	if _, exists := pl.links[page]; !exists {
		pl.links[page] = []string{}
	}
	pl.links[page] = append(pl.links[page], link)
}

// Exists checks if a link exists for the page
func (pl *PageLinks) Exists(page, link string) bool {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	for _, l := range pl.links[page] {
		if l == link {
			return true
		}
	}
	return false
}

// GetLinks gets the links for the page
func (pl *PageLinks) GetLinks(page string) []string {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.links[page]
}


// WIKIRACER
// WikiRacer finds the shortest path between two Wikipedia pages
type WikiRacer struct {
	startURL   string
	endURL     string
	visited    map[string]bool
	queue      []string
	pageLinks  *PageLinks
	pathToLink map[string]string
}

// NewWikiRacer creates a new WikiRacer instance
func NewWikiRacer(startURL, endURL string) *WikiRacer {
	return &WikiRacer{
		startURL:   startURL,
		endURL:     endURL,
		visited:    make(map[string]bool),
		queue:      []string{startURL},
		pageLinks:  NewPageLinks(),
		pathToLink: make(map[string]string),
	}
}

// FindShortestPath starts the BFS to find the shortest path
func (wr *WikiRacer) FindShortestPath() ([]string, error) {
	for len(wr.queue) > 0 {
		currentPage := wr.queue[0]
		wr.queue = wr.queue[1:]

		if currentPage == wr.endURL {
			return wr.buildPath(), nil
		}

		links, err := wr.fetchLinks(currentPage)
		if err != nil {
			return nil, err
		}

		for _, link := range links {
			if !wr.visited[link] {
				wr.visited[link] = true
				wr.queue = append(wr.queue, link)
				wr.pathToLink[link] = currentPage
				if link == wr.endURL {
					return wr.buildPath(), nil
				}
			}
		}
	}
	return nil, fmt.Errorf("no path found from %s to %s", wr.startURL, wr.endURL)
}

// fetchLinks retrieves the distinct Wikipedia links from the given page
func (wr *WikiRacer) fetchLinks(pageURL string) ([]string, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fetching page %s failed with status: %d", pageURL, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	doc.Find("p a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
			link := "https://en.wikipedia.org" + href
			if !wr.pageLinks.Exists(pageURL, link) {
				wr.pageLinks.Add(pageURL, link)
				links = append(links, link)
				}
		}
	})

	return links, nil
}

// buildPath reconstructs the path from the start URL to the end URL
func (wr *WikiRacer) buildPath() []string {
	path := []string{}
	for link := wr.endURL; link != ""; link = wr.pathToLink[link] {
		path = append([]string{link}, path...)
	}
	return path
}

// main function to simulate input and execute the WikiRacer
func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: mywikiracer <startURL> <endURL>")
        os.Exit(1)
    }
    start := os.Args[1]
    end := os.Args[2]

    racer := NewWikiRacer(start, end)
    startTime := time.Now()

    path, err := racer.FindShortestPath()
    if err != nil {
        log.Fatalf("Error finding path: %v", err)
    }

    duration := time.Since(startTime)
    fmt.Printf("Path found: %v\n", path)
    fmt.Printf("Time taken: %v\n", duration)
}