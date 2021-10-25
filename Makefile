# Go parameters
BINARY_NAME=go-rest-template
GO_CMD=go
GOFMT_CMD=gofmt
MOCKGEN_CMD=mockery
STATICCHECK_CMD=staticcheck

all: test build
install:
	$(GO_CMD) install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.8.3
	$(GO_CMD) install github.com/vektra/mockery/v2@latest
	$(GO_CMD) install honnef.co/go/tools/cmd/staticcheck@latest
build:
	$(GO_CMD) build -o $(BINARY_NAME)
generate:
	oapi-codegen --generate 'types' --package 'api' api/openapi.yaml > api/openapi.go
	$(MOCKGEN_CMD) --dir internal/app --output internal/app/mocks --all
	$(MOCKGEN_CMD) --dir internal/redis --output internal/redis/mocks --all
	$(MOCKGEN_CMD) --dir internal/task --output internal/task/mocks --all
format:
	$(GOFMT_CMD) -w -s .
check:
	$(GO_CMD) vet ./...
	$(STATICCHECK_CMD) ./...
test:
	$(GO_CMD) test -race -covermode=atomic -coverprofile cover.out ./...
testshort:
	$(GO_CMD) test -short -race -covermode=atomic -coverprofile cover.out ./...
coverreport:
	$(GO_CMD) tool cover -html=cover.out -o cover.html
clean:
	$(GO_CMD) clean
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)
