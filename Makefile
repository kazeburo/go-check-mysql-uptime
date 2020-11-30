VERSION=0.0.6
LDFLAGS=-ldflags "-X main.version=${VERSION}"
GO111MODULE=on

all: check-mysql-uptime

.PHONY: check-mysql-uptime

check-mysql-uptime: check-mysql-uptime.go
	go build $(LDFLAGS) -o check-mysql-uptime

linux: check-mysql-uptime.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o check-mysql-uptime

check:
	go test ./...

fmt:
	go fmt ./...

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin master
