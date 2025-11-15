package urlmanager_test

import (
	"cmd/crawler/internal/urlmanager"
	"testing"
)

func TestUrlManager_ShouldVisit(t *testing.T) {
	urlManager, err := urlmanager.New("https://crawlme.monzo.com/")
	if err != nil {
		t.Fatalf("Failed to create UrlManager: %v", err)
	}

	shouldVisit, normalizedUrl, err := urlManager.ShouldVisit("https://crawlme.monzo.com/about")
	if err != nil {
		t.Fatalf("Failed to check if URL should be visited: %v", err)
	}
	if !shouldVisit {
		t.Fatalf("URL should be visited")
	}
	if normalizedUrl != "https://crawlme.monzo.com/about.html" {
		t.Fatalf("Normalized URL is incorrect")
	}
}
func TestUrlManager_MarkAsVisited(t *testing.T) {

	urlManager, err := urlmanager.New("https://crawlme.monzo.com/")
	if err != nil {
		t.Fatalf("Failed to create UrlManager: %v", err)
	}

	url := "https://crawlme.monzo.com/about"

	// Before marking: should be visitable
	shouldVisit, _, err := urlManager.ShouldVisit(url)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !shouldVisit {
		t.Error("URL should be visitable before marking")
	}

	// Mark as visited
	err = urlManager.MarkAsVisited(url)
	if err != nil {
		t.Fatalf("Failed to mark URL as visited: %v", err)
	}

	// After marking: should NOT be visitable
	shouldVisit, _, err = urlManager.ShouldVisit(url)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if shouldVisit {
		t.Error("URL should NOT be visitable after marking as visited")
	}
}
