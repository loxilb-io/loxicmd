.DEFAULT_GOAL := build
bin=loxicmd

build: 
	@go build -o ${bin}

test: 
	go test

check:
	go test

run:
	./$(bin)

