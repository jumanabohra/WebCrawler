package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"cmd/crawler/internal/crawler"
)

func main() {
	// Parse command line arguments
	url := flag.String("url", "", "URL to crawl")
	workers := flag.Int("workers", 10, "Maximum number of concurrent workers")
	timeout := flag.Int("timeout", 30, "HTTP request timeout in seconds")
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	memProfile := flag.String("memprofile", "", "write memory profile to file")
	flag.Parse()

	// CPU profiling
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatalf("Failed to create CPU profile file: %v", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("Failed to start CPU profile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Memory profiling (deferred until end of execution)
	if *memProfile != "" {
		defer func() {
			f, err := os.Create(*memProfile)
			if err != nil {
				log.Fatalf("Failed to create memory profile file: %v", err)
			}
			defer f.Close()
			runtime.GC() // Get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatalf("Failed to write memory profile: %v", err)
			}
		}()
	}

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
