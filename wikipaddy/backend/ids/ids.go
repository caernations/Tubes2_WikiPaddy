package ids

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WikiRacerIDS struct {
	startURL      string
	endURL        string
	visited       map[string]int
	pageLinks     map[string][]string
	maxDepth      int
	linksExamined int // Number of links examined
}

func NewWikiRacerIDS(startURL, endURL string) *WikiRacerIDS {
	return &WikiRacerIDS{
		startURL:      startURL,
		endURL:        endURL,
		visited:       make(map[string]int),
		pageLinks:     make(map[string][]string),
		maxDepth:      0,
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

func (wr *WikiRacerIDS) depthLimitedSearch(currentURL string, depth int, path []string) (bool, []string) {
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

	for _, link := range links {
		if prevDepth, visited := wr.visited[link]; !visited || depth+1 < prevDepth {
			found, path := wr.depthLimitedSearch(link, depth+1, path)
			if found {
				// Return the complete path
				return true, path
			}
		}
	}

	return false, nil
}

func (wr *WikiRacerIDS) FindShortestPathUsingIDS() ([]string, error) {
	for {
		found, path := wr.depthLimitedSearch(wr.startURL, 0, []string{})
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
func printPath(path []string){
	if len(path) == 0 {
		return
	}

	fmt.Printf("%s", extractArticleTitle(path[0]))
	for i := 1; i < len(path); i++ {
		fmt.Printf(" -> %s", extractArticleTitle(path[i]))
	}
	fmt.Println()
}