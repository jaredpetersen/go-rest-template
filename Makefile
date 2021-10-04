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
	$(MOCKGEN_CMD) --dir internal/app --output internal/app/mocks --all
	$(MOCKGEN_CMD) --dir internal/redis --output internal/redis/mocks --all
	$(MOCKGEN_CMD) --dir internal/task --output internal/task/mocks --all
test: generate
	$(GO_TEST) -coverprofile cover.out ./...
	$(GO_COVER) -html=cover.out -o cover.html
testshort: generate
	$(GO_TEST) -short -coverprofile cover.out ./...
	$(GO_COVER) -html=cover.out -o cover.html
clean:
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)
