BINARY_NAME=bookkeeping

install:
	go install

clean:
	rm -f bin/*

build:
	go build -o bin/${BINARY_NAME}

dockerImage:
	docker build -t myth-bookkeeping .

dockerDevImage:
	docker build -t myth-bookkeeping . -f dev.dockerfile
