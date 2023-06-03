.ONESHELL: # Applies to every targets in the file!

BINARY_NAME=parser
DATASIZE:=100000

.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

build: clean # Build the binary
	@cd parser
	@go build -o ${BINARY_NAME} main.go

run: build # Run the parser using ../data/data.txt as input
	@echo "data/data.txt" | ./parser/${BINARY_NAME}

clean: # Clean the project
	@cd parser
	@go clean -testcache
	@rm -f ${BINARY_NAME}

test: build # Run Tests
	@cd parser
	@go test ./...

dep: # Download deps
	@cd parser
	@go mod download

vet: # Run Go Bet
	@cd parser
	@go vet

.PHONY: data
data: # Regenerate the data file - depending on the size this might take a while
	cd data
	rm -f data.txt
	seq ${DATASIZE} | xargs -P 10 -I {} sh -c "echo http://api.tech.com/item/{} \$$(od -N 4 -t uL -An /dev/urandom | tr -d \" \") >> data.txt"