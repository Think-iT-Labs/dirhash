.PHONY: build-cli testacc

build-cli:
	go build -o bin ./cmd

build-provider:
	go build -o bin

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
