PHONY:build

build-race:
	go build -race -o ./bin/myhttp

build: clean
	go build -o ./bin/myhttp

PHONY:clean
clean:
	rm -f ./bin/*

PHONY:test
test:
	go clean -testcache
	go test ./internal/...

PHONY:check
check:
	go vet ./...

PHONY:install
install:
	go install ./...