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
		URL:     "https://www.crawlme.monzo.com/",
		Workers: 2, // Use fewer workers for testing
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create crawler: %v", err)
	}

	// TODO: To use mockReporter, Crawler.New() would need to accept
	// an optional Reporter parameter, or add a SetReporter() method.
	// For now, this test verifies the crawler runs successfully.
	//
	// Example of how it would work with dependency injection:
	// c, err := crawler.NewWithReporter(config, mockReporter)
	// Then after Crawl():
	// mockReporter.Wait()
	// results := mockReporter.GetResults()
	// if len(results) == 0 {
	//     t.Error("Expected at least one result")
	// }

	err = c.Crawl()
	if err != nil {
		t.Fatalf("Failed to crawl: %v", err)
	}

	// Integration test verifies the crawler runs to completion
	// The mockReporter is ready to use once dependency injection is added
	_ = mockReporter // Suppress unused variable warning
}
