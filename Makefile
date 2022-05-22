test:
	go test -v -race ./...

fetch:
	CLDR_VERSION=40 ./fetch-cldr.sh

gen: 
	CLDR_VERSION=40 go generate

gen-fetch: fetch gen

clean:
	rm ./pkg/tables.go

clean-all: clean
	rm ./core.zip