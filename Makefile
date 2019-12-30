.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./bin
	
build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/web com.napicella/hello-world/web
	GOOS=linux GOARCH=amd64 go build -o ./bin/authorizer com.napicella/hello-world/authorizer
	GOOS=linux GOARCH=amd64 go build -o ./bin/client com.napicella/hello-world/web/client