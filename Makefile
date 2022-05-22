test:
	go test -v -race ./...

fetch:
	CLDR_VERSION=40 ./fetch-cldr.sh

gen: fetch
	CLDR_VERSION=40 go generate

clean:
	rm ./core.zip

clean-all: clean
	rm ./pkg/tables.go