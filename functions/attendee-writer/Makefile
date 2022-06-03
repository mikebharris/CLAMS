export SHELL:=/bin/bash

.PHONY: build
build:
	GOOS=linux go build -o main .

.PHONY: target
target:
	GOARCH=amd64 GOOS=linux go build -o main .

.PHONY: int-test
int-test: build
	cp ./main ./service-tests/
	cd service-tests && go test

.PHONY: unit-test
unit-test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o ./coverage.html

.PHONY: test
test: unit-test int-test