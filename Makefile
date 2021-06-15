.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/createBreach createBreach/main.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose
