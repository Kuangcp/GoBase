BINARY_NAME=mybook

install:
	go install

clean:
	rm -f bin/*

build:
	go build -o bin/${BINARY_NAME}

run:
	go build -o bin/${BINARY_NAME} && ./bin/${BINARY_NAME} -s

dockerImage:
	docker build -t mybook .

dockerDevImage:
	docker build -t mybook . -f dev.dockerfile
