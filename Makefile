.PHONY: all clean ci test lint format checkformat

LINUX_AMD64=GOOS=linux GOARCH=amd64
LINUX_ARM64=GOOS=linux GOARCH=arm
DARWIN_AMD64=GOOS=darwin GOARCH=amd64
LDFLAGS="-s -w"

all: bin/jqrp.linux-amd64 bin/jqrp.linux-arm64 bin/jqrp.darwin-amd64

bin/jqrp.linux-amd64: bin */**.go
	env $(LINUX_AMD64) go build -ldflags=$(LDFLAGS) -o $@ cmd/jqrp.go

bin/jqrp.linux-arm64: bin */**.go
	env $(LINUX_ARM64) go build -ldflags=$(LDFLAGS) -o $@ cmd/jqrp.go

bin/jqrp.darwin-amd64: bin */**.go
	env $(DARWIN_AMD64) go build -ldflags=$(LDFLAGS) -o $@ cmd/jqrp.go

bin:
	mkdir $@

clean:
	rm -rf bin

ci: test lint

test:
	go test -v ./...

lint:
	golint -set_exit_status ./...

format:
	go fmt ./...

checkformat:
	test -z $(shell gofmt -l .)
