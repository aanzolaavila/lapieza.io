git-commit := $(shell git rev-list -1 HEAD)
ldflags := "-s -w -X main.GitCommit=${git-commit}"
gcflags := -G=3
flags := -ldflags=${ldflags} -gcflags=${gcflags}

.PHONY: build
build: bin clean vendor fmt
	go build ${flags} -o bin cmd/run.go

bin:
	mkdir -p bin

.PHONY: docker-lint
docker-lint:
	docker run --rm -i ghcr.io/hadolint/hadolint < Dockerfile

.PHONY: fmt
fmt: staticcheck
	go fmt ./...

.PHONY: staticcheck
staticcheck: vet
	go install honnef.co/go/tools/cmd/staticcheck@latest
	$(GOPATH)/bin/staticcheck ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: vendor
vendor: tidy
	go mod vendor

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf bin/*
