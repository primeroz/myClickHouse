BINARY_NAME=parser

.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

build: # Build the binary
	go build -o ${BINARY_NAME} main.go

run: build # Run the parser using ../data/output.txt as input
	echo "../data/output.txt" | ./${BINARY_NAME}

clean: # Clean the project
	go clean -testcache
	rm -f ${BINARY_NAME}

test: build # Run Tests
	go test ./...

dep: # Download deps
	go mod download

vet: # Run Go Bet
	go vet

lint: # Run Linting
	golangci-lint run --enable-all