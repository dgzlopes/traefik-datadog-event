.PHONY: lint vendor test

export GO111MODULE=on

lint:
	golangci-lint run

vendor:
	go mod vendor

test:
	go test -v datadogevent_test.go datadogevent.go