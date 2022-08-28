.PHONY: build-cli testacc

build-cli:
	go build -o bin/dirhash ./cmd

build-provider:
	go build -o bin/dirhash-terraform-provider

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
