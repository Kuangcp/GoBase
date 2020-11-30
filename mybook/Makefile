BINARY_NAME=mybook.bin

clean:
	rm -f bin/${BINARY_NAME}

run:
	sh run.sh

runS:
	sh run.sh -s

buildAlpine:
	CGO_ENABLED=0 GOOS=linux go build -o bin/${BINARY_NAME}

install:
	mkdir data && echo -e 'db: \n    file: ./data/main.db\ndebug: true' >> data/mybook.yml && go test -v -test.run TestInit

dockerImage:
	CGO_ENABLED=0 GOOS=linux go build -o bin/${BINARY_NAME} && docker build -t mybook .

dockerDevImage:
	docker build -t mybook . -f dev.dockerfile
