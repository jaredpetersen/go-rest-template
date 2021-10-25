# Go parameters
BINARY_NAME=go-rest-template
GO_CMD=go
GO_INSTALL=$(GO_CMD) install
GO_BUILD=$(GO_CMD) build
GO_GEN=$(GO_CMD) generate
GO_FMT=gofmt
GO_TEST=$(GO_CMD) test
GO_COVER=$(GO_CMD) tool cover
GO_CLEAN=$(GO_CMD) clean
MOCKGEN_CMD=mockery

all: test build
install:
	$(GO_INSTALL) github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.8.3
build:
	$(GO_BUILD) -o $(BINARY_NAME)
generate:
	oapi-codegen --generate 'types' --package 'api' api/openapi.yaml > api/openapi.go
	$(MOCKGEN_CMD) --dir internal/app --output internal/app/mocks --all
	$(MOCKGEN_CMD) --dir internal/redis --output internal/redis/mocks --all
	$(MOCKGEN_CMD) --dir internal/task --output internal/task/mocks --all
format:
	$(GO_FMT) -w -s .
test:
	$(GO_TEST) -coverprofile cover.out ./...
	$(GO_COVER) -html=cover.out -o cover.html
testshort:
	$(GO_TEST) -short -coverprofile cover.out ./...
	$(GO_COVER) -html=cover.out -o cover.html
clean:
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)
