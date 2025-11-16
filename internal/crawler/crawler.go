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
	"sync/atomic"
	"time"
)

type (
	Fetcher interface {
		Fetch(url string) ([]byte, error)
	}

	Parser interface {
		Parse([]byte) ([]string, error)
	}

	UrlManager interface {
		ShouldVisit(url string, baseURL string) (bool, string, error)
		MarkAsVisited(url string) error
	}

	Reporter interface {
		Start(results chan types.Result)
	}

	Config struct {
		URL     string
		Workers int
		Timeout time.Duration
	}

	Crawler struct {
		config      Config
		fetcher     Fetcher
		parser      Parser
		urlmanager  UrlManager
		reporter    Reporter
		urlQueue    chan string
		resultQueue chan types.Result
		wg          sync.WaitGroup
		activeWork  int64 // Counter for URLs currently being processed (atomic)
	}
)

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
		urlQueue:    make(chan string, 1000),       // Buffered to prevent blocking
		resultQueue: make(chan types.Result, 1000), // Buffered to prevent blocking
	}, nil
}

// SetReporter sets a custom reporter for testing purposes
func (c *Crawler) SetReporter(r Reporter) {
	c.reporter = r
}

// SetFetcher sets a custom fetcher for testing purposes
func (c *Crawler) SetFetcher(f Fetcher) {
	c.fetcher = f
}

// Crawl the URL
func (c *Crawler) Crawl() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start completion detector
	c.wg.Add(1)
	go c.detectCompletion(ctx)

	// Start worker goroutines
	for i := 0; i < c.config.Workers; i++ {
		c.wg.Add(1)
		go c.worker(ctx)
	}

	// Start the initial URL
	atomic.AddInt64(&c.activeWork, 1) // Track initial URL
	c.urlQueue <- c.config.URL

	// Start reporter first so it can read results as they come in
	go c.reporter.Start(c.resultQueue)

	// Wait for all workers to finish
	c.wg.Wait()

	// Close resultQueue to signal reporter that no more results are coming
	close(c.resultQueue)

	return nil
}

// detectCompletion monitors when all work is done using a counter-based approach
func (c *Crawler) detectCompletion(ctx context.Context) {
	defer c.wg.Done()

	// Ticker to periodically check if we're done
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Check if all work is done: counter is 0 and queue is empty
			if atomic.LoadInt64(&c.activeWork) == 0 && len(c.urlQueue) == 0 {
				// All work is done, close urlQueue to signal workers to finish
				close(c.urlQueue)
				return
			}
		}
	}
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

			// Process the URL (decrement counter when done, regardless of outcome)
			c.processURL(url)
			atomic.AddInt64(&c.activeWork, -1)
		}
	}
}

// ProcessURL processes a single URL
func (c *Crawler) processURL(url string) {

	// Check if the URL should be visited (no base URL needed for initial URL)
	shouldVisit, normalizedUrl, err := c.urlmanager.ShouldVisit(url, "")
	if err != nil {
		log.Fatalf("Failed to check if URL should be visited: %v", err)
		return
	}

	if !shouldVisit {
		return
	}

	// Mark the URL as visited
	if err := c.urlmanager.MarkAsVisited(normalizedUrl); err != nil {
		log.Fatalf("Failed to mark URL as visited: %v", err)
		return
	}

	// Fetch the HTML body
	htmlBody, err := c.fetcher.Fetch(normalizedUrl)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
		return
	}

	// Parse the HTML body for links
	links, err := c.parser.Parse(htmlBody)
	if err != nil {
		log.Fatalf("Failed to parse HTML body: %v", err)
		return
	}

	// Queue the new URLs directly to urlQueue (buffered, so won't block)
	// Pass the current page URL (normalizedUrl) as base URL for resolving relative links
	for _, link := range links {
		if shouldVisit, normalizedLink, _ := c.urlmanager.ShouldVisit(link, normalizedUrl); shouldVisit {
			atomic.AddInt64(&c.activeWork, 1) // Increment counter for each new URL
			c.urlQueue <- normalizedLink
		}
	}

	// Send the result to the result queue
	c.resultQueue <- types.Result{
		URL:   normalizedUrl,
		Links: links,
	}
}
