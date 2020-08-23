BINARY_NAME=mybook

build:
	go build -o bin/${BINARY_NAME}

buildSingle:
	statik -f -src=conf/static/ -dest app/common/ && go build -o bin/${BINARY_NAME}

buildSingleWin:
	statik -f -src=conf/static/ -dest app/common/ && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}.exe
#  -a -ldflags "-linkmode external -extldflags '-static' -s -w"

clean:
	rm -f bin/mybook

run:
	statik -f -src=conf/static/ -dest app/common/ && go build -o bin/${BINARY_NAME} && ./bin/${BINARY_NAME} -s

buildAlpine:
	CGO_ENABLED=0 GOOS=linux go build -o bin/${BINARY_NAME}

install:
	mkdir data && echo -e 'db: \n    file: ./data/main.db\ndebug: true' >> data/mybook.yml && go test -v -test.run TestInit

dockerImage:
	CGO_ENABLED=0 GOOS=linux go build -o bin/${BINARY_NAME} && docker build -t mybook .

dockerDevImage:
	docker build -t mybook . -f dev.dockerfile
