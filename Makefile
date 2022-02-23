PKG = github.com/jeesmon/operator-utils

all: fmt vet tidy generate test

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

generate: deepcopy-gen
	$(DEEPCOPY_GEN) -i $(PKG)/status -h hack/boilerplate.go.txt --output-base . --trim-path-prefix $(PKG) --output-package $(PKG)/status -O zz_deepcopy -v 10

DEEPCOPY_GEN = $(shell pwd)/bin/deepcopy-gen
.PHONY: deepcopy-gen
deepcopy-gen:
	$(call go-get-tool,$(DEEPCOPY_GEN),k8s.io/code-generator/cmd/deepcopy-gen@v0.24.0-alpha.3)

test:
	go test ./... -coverprofile cover.out -v

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
