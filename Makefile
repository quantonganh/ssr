lint:
	golangci-lint run -v ./...

test:
	go test -v ./...

integration-test:
	go test -tags integration -v -run TestScanService ./...

build:
	CGO_ENABLED=0 go build -v -ldflags="-s -w" -o ssr cmd/ssr/main.go