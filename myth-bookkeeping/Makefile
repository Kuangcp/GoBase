BINARY_NAME=myth-bookkeeping

build:
	go build -o bin/${BINARY_NAME}

dockerImage:
	docker build -t myth-bookkeeping .

dockerDevImage:
	docker build -t myth-bookkeeping . -f dev.dockerfile

clean:
	rm bin/${BINARY_NAME}
