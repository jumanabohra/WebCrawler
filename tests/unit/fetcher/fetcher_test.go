package fetcher_test

import (
	"cmd/crawler/internal/fetcher"
	"testing"
	"time"
)

func TestFetcher_Fetch(t *testing.T) {
	fetcher := fetcher.New(10 * time.Second)
	body, err := fetcher.Fetch("https://crawlme.monzo.com/")
	if err != nil {
		t.Fatalf("Failed to fetch URL: %v", err)
	}
	if len(body) == 0 {
		t.Fatalf("Body is empty")
	}
}
