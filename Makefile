GO = GO111MODULE=on go
MODULE = "stock-min-max"

fmt:
	${GO} fmt ./...

vet: fmt
	${GO} vet ./...

clean: vet
	rm -rf ./bin
	${GO} mod tidy

build: clean
	${GO} build -o ./bin/${MODULE} ./...

test: build
	${GO} test -v -cover ./...

run: clean
	${GO} run ./... -p 80

.PHONY: run build clean vet fmt
