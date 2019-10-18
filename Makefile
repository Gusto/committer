DEFAULT_GOAL = committer

.PHONY: committer
committer:
	go build

.PHONY: deps
deps:
	command -V dep >/dev/null 2>&1 || go get -u github.com/golang/dep
	dep ensure

.PHONY: tools
tools:
	command -V golangci-lint >/dev/null 2>&1 || go install -i ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint: tools
	golangci-lint run

.PHONY: test
test:
	go test ./...
