VERSION=0.0.6
LDFLAGS=-ldflags "-X main.version=${VERSION}"
GO111MODULE=on

all: check-spf-and-reserve-lookup

.PHONY: check-spf-and-reserve-lookup


check-spf-and-reserve-lookup: check-spf-and-reserve-lookup.go
	go build $(LDFLAGS) -o check-spf-and-reserve-lookup check-spf-and-reserve-lookup.go

linux: check-spf-and-reserve-lookup.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o check-spf-and-reserve-lookup check-spf-and-reserve-lookup.go

fmt:
	go fmt ./...

clean:
	rm -rf check-spf-and-reserve-lookup

check:
	go test ./...

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin master
