VERSION = base/version.go

branch = $(shell git rev-parse --abbrev-ref HEAD)
revision = $(shell git rev-parse HEAD)
dirty = $(shell test -n "`git diff --shortstat 2> /dev/null | tail -n1`" && echo "*")
base = github.com/zero-os/0-fs
ldflags = '-w -s -X $(base).Branch=$(branch) -X $(base).Revision=$(revision) -X $(base).Dirty=$(dirty) -extldflags "-static"'

build:
	go build -ldflags $(ldflags)
