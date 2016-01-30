VERSION=0.0.1

all: check-mysql-uptime

.PHONY: check-mysql-uptime

gom:
	go get -u github.com/mattn/gom

bundle:
	gom install

check-mysql-uptime: check-mysql-uptime.go
	gom build -o check-mysql-uptime

linux: check-mysql-uptime.go
	GOOS=linux GOARCH=amd64 gom build -o check-mysql-uptime

fmt:
	go fmt ./...

dist:
	git archive --format tgz HEAD -o check-mysql-uptime-$(VERSION).tar.gz --prefix check-mysql-uptime-$(VERSION)/

clean:
	rm -rf check-mysql-uptime check-mysql-uptime-*.tar.gz

