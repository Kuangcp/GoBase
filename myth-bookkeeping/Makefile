BINARY_NAME=myth-bookkeeping

build:
	go build -o bin/${BINARY_NAME}

image:
	docker build -t myth-bookkeeping .

devImage:
	docker build -t myth-bookkeeping . -f dev.dockerfile