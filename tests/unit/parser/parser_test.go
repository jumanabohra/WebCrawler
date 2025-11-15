package parser_test

import (
	"cmd/crawler/internal/parser"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := parser.New()
	body := []byte("<html><body><a href='https://crawlme.monzo.com/'>About</a></body></html>")
	links, err := parser.Parse(body)
	if err != nil {
		t.Fatalf("Failed to parse body: %v", err)
	}
	if len(links) == 0 {
		t.Fatalf("No links found")
	}
	if links[0] != "https://crawlme.monzo.com/" {
		t.Fatalf("Link is incorrect")
	}
}
