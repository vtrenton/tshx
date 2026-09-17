.PHONY: build install fmt vet snapshot clean

build:
	go build -o tshx .

install:
	go install .

fmt:
	gofmt -l .

vet:
	go vet ./...

snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -f tshx
	rm -rf dist
