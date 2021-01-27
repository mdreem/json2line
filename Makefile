VERSION=$(shell git describe --tags --abbrev=0)
COMMIT=$(shell git rev-list -1 HEAD)

build:
	go build -o bin/main main.go

run:
	go run main.go

test:
	go test -v ./... -covermode=count -coverprofile=coverage.out

lint:
	golangci-lint run --config=.github/linters/golangci.yml

clean:
	rm -r bin/**

compile:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'github.com/mdreem/json2line/cmd.Version=$(VERSION)' -X 'github.com/mdreem/json2line/cmd.GitCommit=$(COMMIT)'" -o bin/darwin-amd64/json2line main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/mdreem/json2line/cmd.Version=$(VERSION)' -X 'github.com/mdreem/json2line/cmd.GitCommit=$(COMMIT)'" -o bin/linux-amd64/json2line main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X 'github.com/mdreem/json2line/cmd.Version=$(VERSION)' -X 'github.com/mdreem/json2line/cmd.GitCommit=$(COMMIT)'" -o bin/windows-amd64/json2line.exe main.go
