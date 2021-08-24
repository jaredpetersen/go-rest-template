# Go parameters
BINARY_NAME=go-rest-template
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_GEN=$(GO_CMD) generate
GO_TEST=$(GO_CMD) test
GO_COVER=$(GO_CMD) tool cover
GO_CLEAN=$(GO_CMD) clean
MOCKGEN_CMD=mockery

all: test build
build:
	$(GO_BUILD) -o $(BINARY_NAME)
generate:
	$(MOCKGEN_CMD) --dir internal/redis --output internal/redis/mocks --all
test: generate
	$(GO_TEST) -coverprofile coverage.out ./...
	$(GO_COVER) -html=coverage.out -o coverage.html
clean:
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)
