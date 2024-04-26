package bfs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	// "github.com/gorilla/mux"
    // "github.com/rs/cors"
)

type PageLinks struct {
	mu    sync.Mutex
	links map[string][]string
}

type Request struct {
    StartURL string `json:"start_url"`
    EndURL   string `json:"end_url"`
}

func HandleRequestBFS(w http.ResponseWriter, r *http.Request) {
    startTitle := r.URL.Query().Get("start_title")
    endTitle := r.URL.Query().Get("end_title")

    if startTitle == "" || endTitle == "" {
        http.Error(w, "start_title and end_title are required", http.StatusBadRequest)
        return
    }

    startURL := "https://en.wikipedia.org/wiki/" + startTitle
    endURL := "https://en.wikipedia.org/wiki/" + endTitle

    wikiRacer := NewWikiRacer(startURL, endURL)
    path, waktuEksekusi ,err := wikiRacer.FindShortestPath()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(path)
}

func NewPageLinks() *PageLinks {
	return &PageLinks{
		links: make(map[string][]string),
	}
}

func (pl *PageLinks) Add(page, link string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	if _, exists := pl.links[page]; !exists {
		pl.links[page] = []string{}
	}
	pl.links[page] = append(pl.links[page], link)
}

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

func (pl *PageLinks) GetLinks(page string) []string {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.links[page]
}

type WikiRacer struct {
	startURL      string
	endURL        string
	visited       map[string]bool
	queue         []string
	pageLinks     *PageLinks
	pathToLink    map[string]string
	linksExamined int
}

func NewWikiRacer(startURL, endURL string) *WikiRacer {
	return &WikiRacer{
		startURL:      startURL,
		endURL:        endURL,
		visited:       make(map[string]bool),
		queue:         []string{startURL},
		pageLinks:     NewPageLinks(),
		pathToLink:    make(map[string]string),
		linksExamined: 0,
	}
}

func (wr *WikiRacer) FindShortestPath() ([]string, time.Duration, error) {
	startTime := time.Now()

	for len(wr.queue) > 0 {
		currentPage := wr.queue[0]
		wr.queue = wr.queue[1:]

		var path []string
		link := currentPage
		for link != wr.startURL {
			path = append([]string{extractArticleTitle(link)}, path...)
			link = wr.pathToLink[link]
		}
		path = append([]string{extractArticleTitle(wr.startURL)}, path...)

		// fmt.Println("Path:", formatPath(path))

		if currentPage == wr.endURL {
			return wr.buildPath(), time.Since(startTime), nil
		}

		links, err := wr.fetchLinks(currentPage)
		if err != nil {
			return nil, 0, err
		}

		for _, link := range links {
			if !wr.visited[link] {
				wr.linksExamined++
				wr.visited[link] = true
				wr.queue = append(wr.queue, link)
				wr.pathToLink[link] = currentPage
				if link == wr.endURL {
					return wr.buildPath(), time.Since(startTime), nil
				}
			}
		}
	}

	return nil, 0, fmt.Errorf("no path found from %s to %s", wr.startURL, wr.endURL)
}


func (wr *WikiRacer) fetchLinks(pageURL string) ([]string, error) {
	resp, err := wr.getWithTimeout(pageURL, 30*time.Second) // Set timeout 30 detik
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching page %s failed with status: %d", pageURL, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	linkCh := make(chan string)

	go func() {
		doc.Find("p a[href]").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
				link := "https://en.wikipedia.org" + href
				linkCh <- link
			}
		})
		close(linkCh)
	}()

	for link := range linkCh {
		if !wr.pageLinks.Exists(pageURL, link) {
			wr.pageLinks.Add(pageURL, link)
			links = append(links, link)
		}
	}

	return links, nil
}

func (wr *WikiRacer) getWithTimeout(url string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	return client.Get(url)
}

func (wr *WikiRacer) buildPath() []string {
	var path []string
	currentPage := wr.endURL

	for currentPage != wr.startURL {
		path = append([]string{currentPage}, path...)
		currentPage = wr.pathToLink[currentPage]
	}

	path = append([]string{wr.startURL}, path...)

	return path
}

func extractArticleTitle(url string) string {
	title := strings.TrimPrefix(url, "https://en.wikipedia.org/wiki/")
	index := strings.Index(title, "/")
	if index != -1 {
		title = title[:index]
	}
	return title
}

func (wr *WikiRacer) LinksExamined() int {
	return wr.linksExamined
}

func formatPath(path []string) string {
	var formattedPath strings.Builder
	formattedPath.WriteString("[")
	for i, p := range path {
		formattedPath.WriteString(p)
		if i < len(path)-1 {
			formattedPath.WriteString(" -> ")
		}
	}
	formattedPath.WriteString("]")
	return formattedPath.String()
}