VERSION = base/version.go

branch = $(shell git rev-parse --abbrev-ref HEAD)
revision = $(shell git rev-parse HEAD)
dirty = $(shell test -n "`git diff --shortstat 2> /dev/null | tail -n1`" && echo "*")
base = github.com/zero-os/0-fs

ldflags = '-extldflags "-lrocksdb -llz4 -lstdc++ -lm -lz -lbz2 -lsnappy"'

build:
	go build -o zbundle -ldflags $(ldflags)
