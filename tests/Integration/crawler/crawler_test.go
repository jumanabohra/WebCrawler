package crawler_test

import (
	"cmd/crawler/internal/crawler"
	"cmd/crawler/internal/types"
	"sync"
	"testing"
	"time"
)

// mockReporter collects results for testing
type mockReporter struct {
	results []types.Result
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func (m *mockReporter) Start(results chan types.Result) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for result := range results {
			m.mu.Lock()
			m.results = append(m.results, result)
			m.mu.Unlock()
		}
	}()
}

func (m *mockReporter) GetResults() []types.Result {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.results
}

func (m *mockReporter) Wait() {
	m.wg.Wait()
}

func TestCrawler_Crawl(t *testing.T) {
	// Create mock reporter to collect results
	mockReporter := &mockReporter{
		results: make([]types.Result, 0),
	}

	// Create crawler
	c, err := crawler.New(crawler.Config{
		URL:     "https://crawlme.monzo.com/",
		Workers: 2, // Use fewer workers for testing
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create crawler: %v", err)
	}

	// Inject mock reporter for testing
	c.SetReporter(mockReporter)

	// Run the crawler
	err = c.Crawl()
	if err != nil {
		t.Fatalf("Failed to crawl: %v", err)
	}

	// Wait for reporter to finish processing all results
	mockReporter.Wait()

	// Verify that results were collected
	results := mockReporter.GetResults()
	if len(results) == 0 {
		t.Error("Expected at least one result from crawling")
	}

	// Verify that the initial URL was crawled
	foundInitialURL := false
	for _, result := range results {
		if result.URL == "https://crawlme.monzo.com" {
			foundInitialURL = true
			break
		}
	}
	if !foundInitialURL {
		t.Error("Expected to find the initial URL in results")
	}
}
