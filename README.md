# WebCrawler

A concurrent web crawler written in Go that crawls websites and extracts links.

## Quick Start

```bash
# Build
make build

# Run
./crawler -url https://crawlme.monzo.com/ -workers 10 -timeout 30 (default make run)
```

## Usage

```bash
./crawler -url <URL> [-workers <num>] [-timeout <seconds>] [-cpuprofile <file>] [-memprofile <file>]
```

**Flags:**

- `-url` (required): Starting URL to crawl
- `-workers`: Number of concurrent workers (default: 10)
- `-timeout`: HTTP request timeout in seconds (default: 30)
- `-cpuprofile`: Enable CPU profiling, write to file
- `-memprofile`: Enable memory profiling, write to file

## Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Generate coverage report
make test-coverage
```

## Performance Profiling

```bash
# CPU profiling
make run-pprof-cpu
make pprof-cpu

# Memory profiling
make run-pprof-mem
make pprof-mem
```

## Project Structure

```
cmd/crawler/          - Main application
internal/
  crawler/            - Core crawler logic
  fetcher/             - HTTP fetching
  parser/              - HTML parsing
  urlmanager/          - URL tracking and normalization
  reporter/            - Result reporting
tests/                 - Unit and integration tests
```

## Features

- Concurrent crawling with configurable worker pool
- URL normalization and duplicate detection
- Per-request timeout handling
- Graceful shutdown with context cancellation
- Test coverage and performance profiling support
