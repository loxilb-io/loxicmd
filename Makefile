.DEFAULT_GOAL := build
bin=loxicmd
dock?=loxilb
SHELL := /bin/bash

loxilbid=$(shell docker ps -f name=$(dock) | grep -w $(dock) | cut  -d " "  -f 1 | grep -iv  "CONTAINER")

# VERSION stamps the release version into cmd.Version. The git tag is the source
# of truth for it -- the source tree carries only a dev placeholder -- so this
# defaults to the tag when building from an exact tag checkout (which is how the
# deb packaging build works) and is empty otherwise. Pass it explicitly with
# `make VERSION=x.y.z` where no tag is reachable. When empty no stamp is applied
# and the placeholder in cmd/root.go shows through, so a plain `go build` or a
# branch build still works unchanged.
VERSION?=$(shell git describe --tags --exact-match 2>/dev/null)
# `override` so that a leading v is stripped even when VERSION came from the
# command line, which a plain assignment cannot do.
override VERSION:=$(patsubst v%,%,$(VERSION))
LDFLAGS:=-X 'loxicmd/cmd.BuildInfo=$(shell date '+%Y_%m_%d')-$(shell git branch --show-current)-$(shell git show --pretty=format:%h --no-patch)'
ifneq ($(VERSION),)
LDFLAGS+= -X 'loxicmd/cmd.Version=$(VERSION)'
endif

build:
	@go build -o ${bin} -ldflags="$(LDFLAGS)"

test: 
	go test

check:
	go test

run:
	./$(bin)

install:
	cp loxicmd /usr/local/sbin/
	/usr/local/sbin/loxicmd completion bash > /etc/bash_completion.d/loxi_completion 
	source /etc/bash_completion.d/loxi_completion

docker-cp: build
	docker cp loxicmd $(loxilbid):/usr/local/sbin/loxicmd