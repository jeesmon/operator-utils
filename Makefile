all: fmt vet tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy
