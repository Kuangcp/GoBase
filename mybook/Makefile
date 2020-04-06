BINARY_NAME=mybook

build:
	go build -o bin/${BINARY_NAME}

clean:
	rm -f bin/*

run:
	go build -o bin/${BINARY_NAME} && ./bin/${BINARY_NAME} -s

buildAlpine:
	CGO_ENABLED=0 GOOS=linux go build -o bin/${BINARY_NAME}

install:
	mkdir data && echo -e 'db: \n    file: ./data/main.db\ndebug: true' >> data/mybook.yml && go test -v -test.run TestInit

cui:
	go build -o bin/${BINARY_NAME} && ./bin/${BINARY_NAME} -c

dockerImage:
	CGO_ENABLED=0 GOOS=linux go build -o bin/${BINARY_NAME} && docker build -t mybook .

dockerDevImage:
	docker build -t mybook . -f dev.dockerfile
