build:
	go build -o crawler cmd/crawler/main.go

run:
	./crawler -url https://crawlme.monzo.com/ -workers 10 -timeout 30

# Test commands
test:
	go test ./... -v

test-unit:
	go test ./tests/unit/... -v

test-integration:
	go test ./tests/Integration/... -v

# Test coverage
test-coverage:
	go test ./... -coverprofile=coverage.out -coverpkg=./internal/...,./cmd/...
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | tail -1

# Performance profiling with pprof
build-pprof: build

run-pprof-cpu:
	@echo "Running crawler with CPU profiling..."
	@echo "Profile will be saved to cpu.prof"
	./crawler -url https://crawlme.monzo.com/ -workers 10 -timeout 30 -cpuprofile cpu.prof || true
	@if [ -f cpu.prof ]; then \
		echo "CPU profile generated: cpu.prof"; \
		echo "View with: make pprof-cpu"; \
	fi

run-pprof-mem:
	@echo "Running crawler with memory profiling..."
	@echo "Profile will be saved to mem.prof"
	./crawler -url https://crawlme.monzo.com/ -workers 10 -timeout 30 -memprofile mem.prof || true
	@if [ -f mem.prof ]; then \
		echo "Memory profile generated: mem.prof"; \
		echo "View with: make pprof-mem"; \
	fi

pprof-cpu:
	@if [ ! -f cpu.prof ]; then \
		echo "CPU profile not found. Run 'make run-pprof-cpu' first."; \
		exit 1; \
	fi
	go tool pprof cpu.prof

pprof-mem:
	@if [ ! -f mem.prof ]; then \
		echo "Memory profile not found. Run 'make run-pprof-mem' first."; \
		exit 1; \
	fi
	go tool pprof mem.prof

clean:
	rm -f crawler
	rm -f *.prof 