package crawler

import (
	"cmd/crawler/internal/fetcher"
	"cmd/crawler/internal/parser"
	"cmd/crawler/internal/reporter"
	"cmd/crawler/internal/types"
	"cmd/crawler/internal/urlmanager"
	"context"
	"log"
	"sync"
	"time"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

type Parser interface {
	Parse([]byte) ([]string, error)
}

type UrlManager interface {
	ShouldVisit(url string) (bool, string, error)
	MarkAsVisited(url string) error
}

type Reporter interface {
	Start(results chan types.Result)
}

type Config struct {
	URL     string
	Workers int
	Timeout time.Duration
}

type Crawler struct {
	config      Config
	fetcher     Fetcher
	parser      Parser
	urlmanager  UrlManager
	reporter    Reporter
	urlQueue    chan string
	resultQueue chan types.Result
	wg          sync.WaitGroup
}

// New creates a new Crawler
func New(config Config) (*Crawler, error) {
	urlMgr, err := urlmanager.New(config.URL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		config:      config,
		fetcher:     fetcher.New(config.Timeout),
		parser:      parser.New(),
		urlmanager:  urlMgr,
		reporter:    reporter.New(),
		urlQueue:    make(chan string),
		resultQueue: make(chan types.Result),
	}, nil
}

// Crawl the URL
func (c *Crawler) Crawl() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker goroutines
	for i := 0; i < c.config.Workers; i++ {
		c.wg.Add(1)
		go c.worker(ctx)
	}

	// Start the initial URL
	c.urlQueue <- c.config.URL

	// Wait for all workers to finish
	c.wg.Wait()
	close(c.urlQueue)
	close(c.resultQueue)

	// Start reporter
	go c.reporter.Start(c.resultQueue)

	return nil
}

func (c *Crawler) worker(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-c.urlQueue:
			if !ok {
				return
			}

			// Check if the URL should be visited
			shouldVisit, normalizedUrl, err := c.urlmanager.ShouldVisit(url)
			if err != nil {
				log.Fatalf("Failed to check if URL should be visited: %v", err)
				continue
			}

			if !shouldVisit {
				continue
			}

			// Mark the URL as visited
			if err := c.urlmanager.MarkAsVisited(normalizedUrl); err != nil {
				log.Fatalf("Failed to mark URL as visited: %v", err)
			}

			// Fetch the HTML body
			htmlBody, err := c.fetcher.Fetch(normalizedUrl)
			if err != nil {
				log.Fatalf("Failed to fetch URL: %v", err)
				continue
			}

			// Parse the HTML body for links
			links, err := c.parser.Parse(htmlBody)
			if err != nil {
				log.Fatalf("Failed to parse HTML body: %v", err)
				continue
			}

			// Queue the new URLs
			for _, link := range links {
				c.urlQueue <- link
			}

			// Send the result (convert to reporter.Result)
			c.resultQueue <- types.Result{
				URL:   normalizedUrl,
				Links: links,
			}
		}
	}
}
