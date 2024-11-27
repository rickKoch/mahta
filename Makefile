all: help

version     ?=  0.0.1
target      ?=  mahta
org         ?=  rickKoch
authorname  ?=  Rick Koch
authoremail ?=  kochrick86@gmail.com
license     ?=  MIT
year        ?=  2023
copyright   ?=  Copyright (c) $(year)

.PHONY: compile
compile: ## Compile for the local architecture ⚙
	@echo "Compiling..."
	go build -ldflags "\
	-X 'github.com/$(org)/$(target).Version=$(version)' \
	-X 'github.com/$(org)/$(target).AuthorName=$(authorname)' \
	-X 'github.com/$(org)/$(target).AuthorEmail=$(authoremail)' \
	-X 'github.com/$(org)/$(target).Copyright=$(copyright)' \
	-X 'github.com/$(org)/$(target).License=$(license)' \
	-X 'github.com/$(org)/$(target).Name=$(target)'" \
	-o $(target) cmd/*.go

.PHONY: install
install: ## Install the program to /usr/bin 🎉
	@echo "Installing..."
	sudo cp $(target) /usr/bin/$(target)

test: clean compile install ## 🤓 Run go tests
	@echo "Testing..."
	go test -v ./...

clean: ## Clean your artifacts 🧼
	@echo "Cleaning..."
	rm -rvf release/*

pre-commit: vet lint imports fmt ## 🤓 pre-comit
	@echo "Testing..."
	go test -v ./...

vet:
	go vet -v ./...

lint:
	golangci-lint run

imports:
	goimports -l -w .

fmt:
	go fmt ./...

.PHONY: format
format:  ## 🤔 Show help messages for make targets
	go fmt ./...

.PHONY: release
release: ## Make the binaries for a GitHub release 📦
	mkdir -p release
	GOOS="linux" GOARCH="amd64" go build -ldflags "-X 'github.com/$(org)/$(target).Version=$(version)'" -o release/$(target)-linux-amd64 cmd/*.go
	GOOS="linux" GOARCH="arm" go build -ldflags "-X 'github.com/$(org)/$(target).Version=$(version)'" -o release/$(target)-linux-arm cmd/*.go
	GOOS="linux" GOARCH="arm64" go build -ldflags "-X 'github.com/$(org)/$(target).Version=$(version)'" -o release/$(target)-linux-arm64 cmd/*.go
	GOOS="linux" GOARCH="386" go build -ldflags "-X 'github.com/$(org)/$(target).Version=$(version)'" -o release/$(target)-linux-386 cmd/*.go
	GOOS="darwin" GOARCH="amd64" go build -ldflags "-X 'github.com/$(org)/$(target).Version=$(version)'" -o release/$(target)-darwin-amd64 cmd/*.go

.PHONY: help
help:  ## 🤔 Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
