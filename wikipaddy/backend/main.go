package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"wikipaddy/bfs"
	"wikipaddy/bi_bfs"
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
		duration = time.Since(startTime)
		linksExamined = bfsInstance.LinksExamined()

	case "bi-bfs":
		biBFSInstance := bi_bfs.NewWikiRacer(startURL, endURL)
		startTime := time.Now()
		path, err = biBFSInstance.FindShortestPath()
		duration = time.Since(startTime)
		linksExamined = biBFSInstance.LinksExamined()
	case "ids":
		idsInstance := ids.NewWikiRacerIDS(startURL, endURL)
		startTime := time.Now()
		path, err = idsInstance.FindShortestPathUsingIDS()
		duration = time.Since(startTime)
		linksExamined = idsInstance.LinksExamined()
	default:
		fmt.Println("Unknown algorithm. Please choose 'bfs' or 'ids'.")
		os.Exit(1)
	}

	fmt.Printf("Algorithm Used: %s\n", algorithm)
	if err != nil { //
		log.Fatalf("Error finding path: %v", err)
	} else {
		fmt.Printf("Jumlah artikel yang diperiksa: %d\n", linksExamined)
		fmt.Println("Jumlah artikel yang dilalui: ", len(path)-1)
		fmt.Println("Path route:")
		fmt.Printf("[")
		for i := 0; i < len(path)-1; i++ {
			parts := strings.Split(path[i], "/")
			rute := parts[len(parts)-1]

			if i < len(path)-2 {
				fmt.Printf("%s -> ", rute)
			} else {
				parts2 := strings.Split(path[i+1], "/")
				rute2 := parts2[len(parts2)-1]
				fmt.Printf("%s -> %s", rute, rute2)
			}
		}
		fmt.Printf("]\n")

		fmt.Printf("Time Taken: %v minutes %v ms\n", duration.Minutes(), duration.Milliseconds()%1000)
	}
}

func ArticleTitleToURL(title string) string {
	formattedTitle := strings.ReplaceAll(title, " ", "_")
	return "https://en.wikipedia.org/wiki/" + formattedTitle
}
