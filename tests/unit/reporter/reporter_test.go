package reporter_test

import (
	"cmd/crawler/internal/reporter"
	"cmd/crawler/internal/types"
	"testing"
)

func TestReporter_Start(t *testing.T) {
	reporter := reporter.New()
	resultChan := make(chan types.Result, 1)

	// Start reporter in goroutine (it blocks on channel read)
	go reporter.Start(resultChan)

	// Send a test result
	resultChan <- types.Result{
		URL:   "https://crawlme.monzo.com/",
		Links: []string{"https://monzo.com/about.html", "https://monzo.com/products.html"},
	}

	// Close channel to signal completion
	close(resultChan)
}
