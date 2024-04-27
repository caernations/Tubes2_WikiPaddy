package ids

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	// "github.com/gorilla/mux"
    // "github.com/rs/cors"
)

type WikiRacerIDS struct {
	startURL      string
	endURL        string
	visited       map[string]int
	maxDepth      int
	linksExamined int // Number of links examined
	cache         *sync.Map
}

type Request struct {
    StartURL string `json:"start_url"`
    EndURL   string `json:"end_url"`
}

func HandleRequestIDS(w http.ResponseWriter, r *http.Request) {
    startTitle := r.URL.Query().Get("start_title")
    endTitle := r.URL.Query().Get("end_title")

    if startTitle == "" || endTitle == "" {
        http.Error(w, "start_title and end_title are required", http.StatusBadRequest)
        return
    }

    startURL := "https://en.wikipedia.org/wiki/" + startTitle
    endURL := "https://en.wikipedia.org/wiki/" + endTitle

    wikiRacer := NewWikiRacerIDS(startURL, endURL)
    path, waktuEksekusi, depth ,err := wikiRacer.FindShortestPathUsingIDS()
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

func NewWikiRacerIDS(startURL, endURL string) *WikiRacerIDS {
	return &WikiRacerIDS{
		startURL:      startURL,
		endURL:        endURL,
		visited:       make(map[string]int),
		maxDepth:      0,
		linksExamined: 0,
		cache:         &sync.Map{},
	}
}

func (wr *WikiRacerIDS) fetchLinks(pageURL string) ([]string, error) {
	if links, ok := wr.cache.Load(pageURL); ok {
		if cachedLinks, ok := links.([]string); ok {
			return cachedLinks, nil
		}
	}

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
			links = append(links, baseURL+href)
		}
	})

	wr.cache.Store(pageURL, links)

	return links, nil
}

func (wr *WikiRacerIDS) depthLimitedSearch(currentURL string, depth int, path []string, timeoutCh <-chan time.Time) (bool, []string) {
	if depth > wr.maxDepth {
		return false, nil
	}

	if currentURL == wr.endURL {
		// Append endURL to the path
		path = append(path, currentURL)
		return true, path
	}

	// Mark the current page as visited at the current depth.
	wr.visited[currentURL] = depth
	wr.linksExamined++

	// Append currentURL to the path
	path = append(path, currentURL)

	// Print the path from startURL to currentURL
	printPath(path)

	links, err := wr.fetchLinks(currentURL)
	if err != nil {
		log.Printf("Failed to fetch links from %s: %v", currentURL, err)
		return false, nil
	}

	// Check if endURL is among the newly fetched links
	for _, link := range links {
		wr.linksExamined++
		if link == wr.endURL {
			// Append endURL to the path and return
			path = append(path, wr.endURL)
			return true, path
		}
	}

	for _, link := range links {
		if prevDepth, visited := wr.visited[link]; !visited || depth+1 < prevDepth {
			found, path := wr.depthLimitedSearch(link, depth+1, path, timeoutCh)
			if found {
				return true, path
			}
		}
	}

	select {
	case <-timeoutCh:
		return false, nil
	default:
		return false, nil
	}
}

func (wr *WikiRacerIDS) FindShortestPathUsingIDS() ([]string, float64, int,error) {
	startTime := time.Now() // Waktu awal eksekusi

	timeoutCh := make(chan time.Time)

	for {
		found, path := wr.depthLimitedSearch(wr.startURL, 0, []string{}, timeoutCh)
		if found {
			return path, time.Since(startTime).Seconds(),len(path), nil // Kembalikan jalur dan waktu eksekusi
		}
		wr.maxDepth++
		wr.visited = make(map[string]int) // Reset visited for the next iteration
	}
}

func (wr *WikiRacerIDS) LinksExamined() int {
	return wr.linksExamined
}

// extractArticleTitle extracts the title of the Wikipedia article from the URL
func extractArticleTitle(url string) string {
	// Remove "https://en.wikipedia.org/wiki/" from the beginning of the URL
	title := strings.TrimPrefix(url, "https://en.wikipedia.org/wiki/")

	// Find the index of the first "/" after "/wiki/"
	index := strings.Index(title, "/")
	if index != -1 {
		// If found "/", take the part before it as the article title
		title = title[:index]
	}

	return title
}

// printPath prints the path from startURL to currentURL
func printPath(path []string) {
	if len(path) == 0 {
		return
	}

	fmt.Printf("%s", extractArticleTitle(path[0]))
	for i := 1; i < len(path); i++ {
		fmt.Printf(" -> %s", extractArticleTitle(path[i]))
	}
	fmt.Println()
}