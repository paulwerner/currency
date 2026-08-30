CLDR_VERSION ?= 40

.PHONY: all build test cover vet fmt fmt-check fetch gen gen-fetch clean clean-all

all: build test

build:
	go build ./...

test:
	go test -v -race ./...

cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then \
		echo "not gofmt-formatted:" $$out; exit 1; \
	fi

fetch:
	CLDR_VERSION=$(CLDR_VERSION) ./fetch-cldr.sh

gen:
	go generate

gen-fetch: fetch gen

clean:
	rm -f coverage.out

clean-all: clean
	rm -f core.zip
