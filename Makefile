.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/addAccountsToBreach addAccountsToBreach/main.go
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/createBreach createBreach/main.go
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/findAccount findAccount/main.go
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/notifyMe notifyMe/main.go
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/notifySubscribersOfBreach notifySubscribersOfBreach/main.go
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/sendSubscriptionEmail sendSubscriptionEmail/main.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose
