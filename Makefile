VERSION=0.0.5
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

all: check-mysql-uptime

.PHONY: check-mysql-uptime

bundle:
	dep ensure

update:
	dep ensure -update

check-mysql-uptime: check-mysql-uptime.go
	go build $(LDFLAGS) -o check-mysql-uptime

linux: main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o check-mysql-uptime

check:
	go test ./...

fmt:
	go fmt ./...

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin master
	goreleaser --rm-dist
