package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"wikipaddy/bfs"
	"wikipaddy/ids"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: <program> <algorithm> <startArticle> <endArticle>")
		os.Exit(1)
	}

	algorithm := os.Args[1]
	startArticle := os.Args[2]
	endArticle := os.Args[3]

	startURL := ArticleTitleToURL(startArticle)
	endURL := ArticleTitleToURL(endArticle)

	var path []string
	var duration time.Duration
	var err error
	var linksExamined int

	switch algorithm {
	case "bfs":
		bfsInstance := bfs.NewWikiRacer(startURL, endURL)
		startTime := time.Now()
		path, err = bfsInstance.FindShortestPath()
		linksExamined = bfsInstance.LinksExamined()
		duration = time.Since(startTime)
	case "ids":
		idsInstance := ids.NewWikiRacerIDS(startURL, endURL)
		startTime := time.Now()
		path, err = idsInstance.FindShortestPathUsingIDS()
		linksExamined = idsInstance.LinksExamined()
		duration = time.Since(startTime)
	default:
		fmt.Println("Unknown algorithm. Please choose 'bfs' or 'ids'.")
		os.Exit(1)
	}

	fmt.Printf("Algorithm Used: %s\n", algorithm)
	if err != nil {
		log.Fatalf("Error finding path: %v", err)
	} else {
		fmt.Printf("Number of links examined: %d\n", linksExamined)
		fmt.Println("Path route:")
		fmt.Printf("[")
		for i := 0; i < len(path)-1; i++ {
			if i < len(path)-2 {
				fmt.Printf("%s -> ", path[i])
			} else {
				fmt.Printf("%s -> %s", path[i], path[i+1])
			}
		}
		fmt.Printf("]\n")
		
		fmt.Printf("Time Taken (ms): %v\n", duration.Milliseconds())
	}
}

func ArticleTitleToURL(title string) string {
	formattedTitle := strings.ReplaceAll(title, " ", "_")
	return "https://en.wikipedia.org/wiki/" + formattedTitle
}
