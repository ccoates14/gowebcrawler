package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strings"
)

var queue = &Queue{}
var targetTerm *string
var startUrl *string
var targetTermLower string
var pagesSeen = 0

type WebPageInfo struct {
	text              string
	linksToOtherPages []string
	pageError         error
}

func main() {
	fmt.Println("Wiki crawl search starting")
	targetTerm = flag.String("targetTerm", "", "The search term for wikipedia")
	startUrl = flag.String("startUrl", "", "The wiki start url")

	flag.Parse()

	fmt.Println("Start url: " + *startUrl)
	fmt.Println("Target term: " + targetTermLower)

	if *targetTerm == "" || *startUrl == "" {
		fmt.Println("Invalid args passed")
		os.Exit(1)
	}

	targetTermLower = strings.ToLower(*targetTerm)

	if targetTermLower == "" {
		os.Exit(1)
	}

	queue.Enqueue(*startUrl)

	crawl()
}

func crawl() {

	for !queue.IsEmpty() {
		//pop url off queue
		url := queue.Dequeue()
		pagesSeen++

		fmt.Println("Searching page: " + url)
		fmt.Println(pagesSeen)

		pageInfo := getWebpageInfo(url)

		if pageInfo.pageError != nil {
			fmt.Println(pageInfo.pageError)
		} else {
			if strings.Contains(pageInfo.text, targetTermLower) {
				fmt.Println("Found search term on url: " + url)
				fmt.Printf("Pages Seen: %d\n", pagesSeen)

				break
			} else {

				for i := 0; i < len(pageInfo.linksToOtherPages); i++ {
					queue.Enqueue(pageInfo.linksToOtherPages[i])
				}
			}
		}
	}
}

func getWebpageInfo(url string) WebPageInfo {
	// Request the HTML page.
	var text string
	var pageError error

	res, err := http.Get(url)
	if err != nil {
		pageError = err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		pageError = errors.New("invalid response")
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		pageError = err
	}

	if pageError != nil {
		return WebPageInfo{
			text:              "",
			linksToOtherPages: []string{},
			pageError:         pageError,
		}
	}

	var urls []string
	// Find the review items
	doc.Find(".mw-content-container").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		s.Find("p").Each(func(i int, s *goquery.Selection) {
			text += s.Text()
		})
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				if strings.HasPrefix(href, "/wiki") {
					urls = append(urls, "https://en.wikipedia.org"+href)
				}
			}
		})
	})

	return WebPageInfo{
		text:              strings.ToLower(text),
		linksToOtherPages: urls,
		pageError:         pageError,
	}
}

/*
	take in a command line arg for target term, and second arg for wiki url
	program will always search wikipedia

	starting with given url crawl wikipedia for search term

	add first page to queue

	create worker thread {
			while queue not empty

				current page = thread safe take url from queue

				does current page contain search term?
					print page url and end program
				else
					find all links on page and add them to the queue

				if thread count < 100
					spawn new worker thread

	}

*/
