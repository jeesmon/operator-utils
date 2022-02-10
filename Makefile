all: fmt vet tidy test

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

test:
	go test ./... -coverprofile cover.out -v
