run.lint:
	clear
	golangci-lint run ./...
	gocritic check ./...
	gocyclo -top 5
mod:
	clear
	go mod tidy
	go mod vendor
run.test:
	go test ./...