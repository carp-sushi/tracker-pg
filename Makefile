.PHONY: all
all: fmt test build

.PHONY: fmt
fmt:
	@gofmt -w ./**/*.go

.PHONY: build
build:
	@go build

.PHONY: test
test:
	@go test -v ./database/repo

.PHONY: clean
clean:
	@go clean

.PHONY: run
run:
	@go run main.go
