BINARY_NAME=mybook

build:
	go build -o bin/${BINARY_NAME}

clean:
	rm -f bin/*

buildAlpine:
	CGO_ENABLED=0 GOOS=linux go build -o bin/${BINARY_NAME}

run:
	go build -o bin/${BINARY_NAME} && ./bin/${BINARY_NAME} -s

dockerImage:
	CGO_ENABLED=0 GOOS=linux go build -o bin/${BINARY_NAME} && docker build -t mybook .

dockerDevImage:
	docker build -t mybook . -f dev.dockerfile
