package bi_bfs

import (
	"fmt"
	// "log"
	"net/http"
	"strings"
	"sync"
	"encoding/json"
	"time"

	"github.com/PuerkitoBio/goquery"
	// "github.com/gorilla/mux"
    // "github.com/rs/cors"
)

// PageLinks represents a mapping from a page to its links
type PageLinks struct {
	mu    sync.Mutex
	links map[string][]string
}

type Request struct {
    StartURL string `json:"start_url"`
    EndURL   string `json:"end_url"`
}

func HandleRequestBiBFS(w http.ResponseWriter, r *http.Request) {
    startTitle := r.URL.Query().Get("start_title")
    endTitle := r.URL.Query().Get("end_title")

    if startTitle == "" || endTitle == "" {
        http.Error(w, "start_title and end_title are required", http.StatusBadRequest)
        return
    }

    startURL := "https://en.wikipedia.org/wiki/" + startTitle
    endURL := "https://en.wikipedia.org/wiki/" + endTitle

    wikiRacer := NewWikiRacer(startURL, endURL)
    path, waktuEksekusi, depth ,err := wikiRacer.FindShortestPath()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(
		map[string]interface{}{
			"path":            path,
			"waktu_eksekusi":  waktuEksekusi,
			"kedalaman": depth,
		},
	)
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

// WikiRacer represents the bidirectional BFS WikiRacer
type WikiRacer struct {
	startURL       string
	endURL         string
	visitedForward map[string]bool
	visitedBackward map[string]bool
	queueForward   []string
	queueBackward   []string
	pathToLinkForward map[string]string
	pathToLinkBackward map[string]string
	linksExamined   int
	pageLinks      *PageLinks
}

// NewWikiRacer creates a new WikiRacer instance
func NewWikiRacer(startURL, endURL string) *WikiRacer {
	return &WikiRacer{
		startURL:         startURL,
		endURL:           endURL,
		visitedForward:   make(map[string]bool),
		visitedBackward:  make(map[string]bool),
		queueForward:     []string{startURL},
		queueBackward:    []string{endURL},
		pathToLinkForward: make(map[string]string),
		pathToLinkBackward: make(map[string]string),
		linksExamined:    0,
		pageLinks:        NewPageLinks(),
	}
}

// FindShortestPath finds the shortest path between start and end URLs using Bi-BFS
func (wr *WikiRacer) FindShortestPath() ([]string, float64, int,error) {
	startTime := time.Now()

	for len(wr.queueForward) > 0 && len(wr.queueBackward) > 0 {
		// Forward BFS
		currentForward := wr.queueForward[0]
		wr.queueForward = wr.queueForward[1:]
		wr.visitedForward[currentForward] = true


		// Backward BFS
		currentBackward := wr.queueBackward[0]
		wr.queueBackward = wr.queueBackward[1:]
		wr.visitedBackward[currentBackward] = true


		// Check for intersection
		if intersectionNode, found := wr.checkIntersection(); found {
			return wr.buildPath(intersectionNode), time.Since(startTime).Seconds(), len(wr.buildPath(intersectionNode)),nil
		}

		// Expand forward BFS
		linksForward, err := wr.fetchLinks(currentForward)
		if err != nil {
			return nil, 0,0, err
		}
		for _, link := range linksForward {
			if !wr.visitedForward[link] {
				wr.visitedForward[link] = true
				wr.queueForward = append(wr.queueForward, link)
				wr.pathToLinkForward[link] = currentForward
				wr.IncrementLinksExamined()
			}
		}

		// Expand backward BFS
		linksBackward, err := wr.fetchLinks(currentBackward)
		if err != nil {
			return nil, 0,0, err
		}
		for _, link := range linksBackward {
			if !wr.visitedBackward[link] {
				wr.visitedBackward[link] = true
				wr.queueBackward = append(wr.queueBackward, link)
				wr.pathToLinkBackward[link] = currentBackward
				wr.IncrementLinksExamined()
			}
		}
	}

	return nil, 0, 0,fmt.Errorf("no path found from %s to %s", wr.startURL, wr.endURL)
}


// checkIntersection checks if there is an intersection between forward and backward BFS
func (wr *WikiRacer) checkIntersection() (string, bool) {
	for url := range wr.visitedForward {
		if wr.visitedBackward[url] {
			return url, true
		}
	}
	return "", false
}

// buildPath reconstructs the path from start URL to end URL
func (wr *WikiRacer) buildPath(intersectionNode string) []string {
	var path []string
	// Build path from startURL to intersectionNode
	current := intersectionNode
	for current != wr.startURL {
		path = append([]string{current}, path...)
		current = wr.pathToLinkForward[current]
	}
	path = append([]string{wr.startURL}, path...)

	// Build path from intersectionNode to endURL
	current = intersectionNode
	for current != wr.endURL {
		current = wr.pathToLinkBackward[current]
		path = append(path, current)
	}

	return path
}

// fetchLinks retrieves distinct Wikipedia links from the given page
func (wr *WikiRacer) fetchLinks(pageURL string) ([]string, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	baseURL := "https://en.wikipedia.org"
	doc.Find("p a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
			link := baseURL + href
			if !wr.pageLinks.Exists(pageURL, link) {
				wr.pageLinks.Add(pageURL, link)
				links = append(links, link)
			}
		}
	})

	return links, nil
}
func (wr *WikiRacer) IncrementLinksExamined(){
	wr.linksExamined++
}

// LinksExamined returns the number of links examined during the search
func (wr *WikiRacer) LinksExamined() int {
	return wr.linksExamined
}