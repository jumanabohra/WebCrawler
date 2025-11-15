package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"cmd/crawler/internal/crawler"
)

func main() {
	// Parse command line arguments
	url := flag.String("url", "", "URL to crawl")
	workers := flag.Int("Workers", 10, "Maximum number of concurrent workers")
	timeout := flag.Int("timeout", 30, "HTTP request timeout in seconds")
	flag.Parse()

	if *url == "" {
		fmt.Fprintf(os.Stderr, "Error: -url flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Crawler configuration
	config := crawler.Config{
		URL:     *url,
		Workers: *workers,
		Timeout: time.Duration(*timeout) * time.Second,
	}

	// Initialize and run crawler
	c, err := crawler.New(config)
	if err != nil {
		log.Fatalf("Failed to initialize crawler: %v", err)
	}

	if err := c.Crawl(); err != nil {
		log.Fatalf("Crawler error: %v", err)
	}
}
