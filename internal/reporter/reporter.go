package reporter

import (
	"cmd/crawler/internal/types"
	"fmt"
)

type Reporter struct {
}

// New creates a new Reporter
func New() *Reporter {
	return &Reporter{}
}

// Start the reporter (matches interface - takes channel as parameter)
func (r *Reporter) Start(results chan types.Result) {
	for result := range results {
		fmt.Printf("Visited: %s\n", result.URL)
		if len(result.Links) > 0 {
			fmt.Printf("Found %d link(s):\n", len(result.Links))
			for _, link := range result.Links {
				fmt.Printf("  - %s\n", link)
			}
		} else {
			fmt.Printf("No links found\n")
		}
		fmt.Println()
	}
}
