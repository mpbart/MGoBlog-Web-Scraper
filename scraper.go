package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var parseUntil string
var shouldScrape bool
var shouldProcess bool
var parseUntilDate time.Time
var directory = "mgoblog.com/"
var baseLink = "http://www.MGoBlog.com"
var nextPage = "?page="
var page = 0

func init() {
	flag.StringVar(&parseUntil, "date", "2015-07-01", "Date to parse until in the format of <YYYY-MM-DD>")
	flag.BoolVar(&shouldScrape, "scrape", false, "Flag indicating whether to scrape mgoblog")
	flag.BoolVar(&shouldProcess, "process", false, "Flag indicating whether to process scraped articles")
	flag.Parse()
	parseUntilDate, _ = time.Parse("2006-01-02", parseUntil)
}

type RawArticle struct {
	article Article
	reader  io.Reader
}

func (m *RawArticle) Work() WordCounter {
	m.article, _ = ProcessArticle(m.reader)
	for _, tag := range m.article.metadata.Tags {
		if ExcludeTagSet.Contains(tag) {
			return WordCounter{}
		}
	}
	retval := CountWords(m.article)
	return retval
}

func main() {
	// Main entry point
	if shouldScrape {
		for {
			webLink := baseLink + nextPage + strconv.Itoa(page)
			_, response, _ := requestPage(webLink)
			page += 1
			metadata, _ := parseMetadataForArticles(response)
			for _, article := range metadata {
				if !fileExists(directory+article.Title) && parseUntilDate.Before(article.Date) {
					articleContents, _, err := requestPage(article.link)
					if err != nil {
						fmt.Println("Error requesting page ", article.link, err)
						continue
					}
					fmt.Println("Requesting/caching article: ", article.Title)
					cachePage(article.Title, articleContents)
				} else {
					break
				}
			}
		}
	}
	if shouldProcess {
		fmt.Println("Processing...")
		processCachedArticles()
	}
}

func processCachedArticles() {
	myChan := make(chan WordCounter)

	go func() {
		for result := range myChan {
			aggregateResults(result)
		}
	}()

	files, _ := ioutil.ReadDir("mgoblog.com")
	var done sync.WaitGroup
	for _, fileInfo := range files {
		page := retrieveCachedPage("mgoblog.com/" + fileInfo.Name())
		done.Add(1)
		go func(page io.Reader) {
			newthing := RawArticle{reader: page}
			results := newthing.Work()
			myChan <- results
			done.Done()
		}(page)
	}
	done.Wait()
	printTopNResults(50)
}

func retrieveCachedPage(title string) io.Reader {
	content, err := ioutil.ReadFile(title)
	if err != nil {
		fmt.Println("Error reading file", title, err)
	}
	return bytes.NewReader(content)
}

func cachePage(title, content string) {
	title = SanitizeFilename(title)
	f, _ := os.Create(directory + title)
	defer f.Close()
	f.WriteString(content)
	f.Sync()
}

func requestPage(link string) (string, io.Reader, error) {

	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("ERROR: failure to get link %s\n", link)
		return "", nil, err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: failure to read response body: %s\n", err)
		return "", nil, err
	}
	return string(contents), bytes.NewReader(contents), nil
}
