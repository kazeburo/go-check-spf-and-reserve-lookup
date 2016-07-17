VERSION=0.0.4
TARGETS_NOVENDOR=$(shell glide novendor)

all: check-spf-and-reserve-lookup

.PHONY: check-spf-and-reserve-lookup

glide:
	go get -u github.com/Masterminds/glide

bundle:
	glide install

check-spf-and-reserve-lookup: check-spf-and-reserve-lookup.go
	GO15VENDOREXPERIMENT=1 go build check-spf-and-reserve-lookup.go

linux: check-spf-and-reserve-lookup.go
	GOOS=linux GOARCH=amd64 GO15VENDOREXPERIMENT=1 go build check-spf-and-reserve-lookup.go

fmt:
	go fmt ./...

dist:
	git archive --format tgz HEAD -o check-spf-and-reserve-lookup-$(VERSION).tar.gz --prefix check-spf-and-reserve-lookup-$(VERSION)/

clean:
	rm -rf check-spf-and-reserve-lookup check-spf-and-reserve-lookup-*.tar.gz

