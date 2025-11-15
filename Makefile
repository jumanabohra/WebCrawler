build:
	go build -o crawler cmd/crawler/main.go

run:
	./crawler -url https://www.monzo.com -workers 10 -timeout 30

clean:
	rm -f crawler