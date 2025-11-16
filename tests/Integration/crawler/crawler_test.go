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

// mockFetcher returns predefined HTML content based on URL
type mockFetcher struct {
	pages map[string]string
}

func newMockFetcher() *mockFetcher {
	return &mockFetcher{
		pages: map[string]string{
			// Root page with relative links
			"https://crawlme.monzo.com": `
				<html>
					<body>
						<a href="about.html">About</a>
						<a href="blog.html">Blog</a>
						<a href="/products.html">Products</a>
					</body>
				</html>`,
			// Blog page with pagination links (relative URLs - CRITICAL TEST)
			// "2.html" should resolve to "https://crawlme.monzo.com/2.html"
			// when current page is "https://crawlme.monzo.com/blog.html"
			"https://crawlme.monzo.com/blog.html": `
				<html>
					<body>
						<a href="post-1.html">Post 1</a>
						<a href="post-2.html">Post 2</a>
						<a href="2.html">Page 2</a>
						<a href="3.html">Page 3</a>
					</body>
				</html>`,
			// Blog page 2 - should be found via relative link "2.html" from blog.html
			// This tests that "2.html" resolves correctly relative to blog.html
			"https://crawlme.monzo.com/2.html": `
				<html>
					<body>
						<a href="post-3.html">Post 3</a>
						<a href="post-4.html">Post 4</a>
						<a href="blog.html">Back to Blog</a>
					</body>
				</html>`,
			// About page
			"https://crawlme.monzo.com/about.html": `
				<html>
					<body>
						<a href="/index.html">Home</a>
					</body>
				</html>`,
			// Products page
			"https://crawlme.monzo.com/products.html": `
				<html>
					<body>
						<a href="product-1.html">Product 1</a>
					</body>
				</html>`,
		},
	}
}

func (m *mockFetcher) Fetch(url string) ([]byte, error) {
	// Normalize URL for lookup (remove trailing slash)
	normalizedURL := url
	if len(normalizedURL) > 0 && normalizedURL[len(normalizedURL)-1] == '/' {
		normalizedURL = normalizedURL[:len(normalizedURL)-1]
	}

	html, exists := m.pages[normalizedURL]
	if !exists {
		// Return empty page for unknown URLs
		return []byte("<html><body></body></html>"), nil
	}
	return []byte(html), nil
}

func TestCrawler_Crawl(t *testing.T) {
	// Create mock reporter to collect results
	mockReporter := &mockReporter{
		results: make([]types.Result, 0),
	}

	// Create mock fetcher with test data
	mockFetcher := newMockFetcher()

	// Create crawler
	c, err := crawler.New(crawler.Config{
		URL:     "https://crawlme.monzo.com",
		Workers: 2,
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create crawler: %v", err)
	}

	// Inject mocks for testing
	c.SetReporter(mockReporter)
	c.SetFetcher(mockFetcher)

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
		t.Fatal("Expected at least one result from crawling")
	}

	// Build a map of crawled URLs for easier checking
	crawledURLs := make(map[string]bool)
	for _, result := range results {
		crawledURLs[result.URL] = true
	}

	// Verify that the initial URL was crawled
	if !crawledURLs["https://crawlme.monzo.com"] {
		t.Error("Expected to find the initial URL in results")
	}

	// Verify that blog.html was crawled (from root page)
	if !crawledURLs["https://crawlme.monzo.com/blog.html"] {
		t.Error("Expected to find blog.html in results")
	}

	// CRITICAL TEST: Verify that 2.html was crawled
	// This tests that relative URLs (like "2.html" from blog.html) are properly resolved
	// When blog.html contains <a href="2.html">, it should resolve to:
	// https://crawlme.monzo.com/2.html (relative to blog.html's location)
	// This would fail if baseURL resolution wasn't working correctly
	if !crawledURLs["https://crawlme.monzo.com/2.html"] {
		t.Error("Expected to find 2.html in results - this tests relative URL resolution with baseURL. " +
			"The link '2.html' from blog.html should resolve to https://crawlme.monzo.com/2.html")
	}

	// Verify that about.html was crawled
	if !crawledURLs["https://crawlme.monzo.com/about.html"] {
		t.Error("Expected to find about.html in results")
	}

	// Verify that products.html was crawled
	if !crawledURLs["https://crawlme.monzo.com/products.html"] {
		t.Error("Expected to find products.html in results")
	}

	t.Logf("Successfully crawled %d pages", len(results))
}
