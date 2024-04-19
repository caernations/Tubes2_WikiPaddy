package ids

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WikiRacerIDS struct {
	startURL   string
	endURL     string
	visited    map[string]int
	pageLinks  map[string][]string
	maxDepth   int
	linksExamined int // Number of links examined
}

func NewWikiRacerIDS(startURL, endURL string) *WikiRacerIDS {
	return &WikiRacerIDS{
		startURL:  startURL,
		endURL:    endURL,
		visited:   make(map[string]int),
		pageLinks: make(map[string][]string),
		maxDepth:  0,
		linksExamined: 1,
	}
}

func (wr *WikiRacerIDS) fetchLinks(pageURL string) ([]string, error) {
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
	return links, nil
}

func (wr *WikiRacerIDS) depthLimitedSearch(currentURL string, depth int) (bool, []string) {
	if depth > wr.maxDepth {
		return false, nil
	}

	if currentURL == wr.endURL {
		return true, []string{wr.endURL}
	}

	// Mark the current page as visited at the current depth.
	wr.visited[currentURL] = depth
	wr.linksExamined++

	links, err := wr.fetchLinks(currentURL)
	if err != nil {
		log.Printf("Failed to fetch links from %s: %v", currentURL, err)
		return false, nil
	}

	for _, link := range links {
		if prevDepth, visited := wr.visited[link]; !visited || depth+1 < prevDepth {
			found, path := wr.depthLimitedSearch(link, depth+1)
			if found {
				// Prepend currentURL to the path and return.
				return true, append([]string{currentURL}, path...)
			}
		}
	}

	return false, nil
}

func (wr *WikiRacerIDS) FindShortestPathUsingIDS() ([]string, error) {
	for {
		found, path := wr.depthLimitedSearch(wr.startURL, 0)
		if found {
			return path, nil
		}
		wr.maxDepth++
		wr.visited = make(map[string]int) // Reset visited for the next iteration
	}
}

func (wr *WikiRacerIDS) LinksExamined() int {
    return wr.linksExamined
}