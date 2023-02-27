.DEFAULT_GOAL := build
bin=loxicmd
dock?=loxilb
SHELL := /bin/bash

loxilbid=$(shell docker ps -f name=$(dock) | grep -w $(dock) | cut  -d " "  -f 1 | grep -iv  "CONTAINER")

build: 
	@go build -o ${bin}

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