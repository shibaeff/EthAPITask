tests:
	go test -race ./internal/...
coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./internal/...
lint:
	golangci-lint run --fix
run:
	go run cmd/eth_validator_api/main.go --config config.yaml server
