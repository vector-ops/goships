build:
	go build -o bin/goships .

run: build
	./bin/goships

test:
	go test -v ./...

clean:
	rm -f bin/goships

clean_logs:
	rm -f logs/*

clean_all: clean clean_logs
