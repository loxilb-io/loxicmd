.DEFAULT_GOAL := build
bin=loxicmd
dock?=loxilb
SHELL := /bin/bash

loxilbid=$(shell docker ps -f name=$(dock) | grep -w $(dock) | cut  -d " "  -f 1 | grep -iv  "CONTAINER")

build: 
	@go build -o ${bin} -ldflags="-X 'loxicmd/cmd.BuildInfo=${shell date '+%Y_%m_%d'}-${shell git branch --show-current}-$(shell git show --pretty=format:%h --no-patch)' -X 'loxicmd/cmd.Version=${shell git describe --tags --abbrev=0}'"

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