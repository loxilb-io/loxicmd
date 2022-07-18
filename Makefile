.DEFAULT_GOAL := build
bin=loxicmd

loxilbid=$(shell docker ps -f name=loxilb | cut  -d " "  -f 1 | grep -iv  "CONTAINER")

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

docker-cp: build
	docker cp loxicmd $(loxilbid):/usr/local/sbin/loxicmd