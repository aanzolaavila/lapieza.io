ldflags := "-s -w"
gcflags := -G=3
flags := -ldflags=${ldflags} -gcflags=${gcflags}

.PHONY: build
build: bin clean vendor fmt build-run build-random

.PHONY: build-run
build-run:
	go build ${flags} -o bin cmd/run/run.go

.PHONY: build-random
build-random:
	go build ${flags} -o bin cmd/random/random.go

bin:
	mkdir -p bin

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
