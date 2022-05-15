all: test build
install:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.10.1
	brew install mockery && brew upgrade mockery
	go install honnef.co/go/tools/cmd/staticcheck@latest
build:
	go build -o go-rest-template
generate:
	# Generate code from OpenAPI schema
	oapi-codegen -generate 'types' -package 'api' api/openapi.yaml > api/openapi.go
	# Generate testing mocks
	mockery --dir internal/app --output internal/app/mocks --all
	mockery --dir internal/redis --output internal/redis/mocks --all
	mockery --dir internal/task --output internal/task/mocks --all
format:
	gofmt -w -s .
check:
	go vet ./...
	staticcheck ./...
test:
	go test -race -timeout 1m -covermode=atomic -coverprofile cover.out ./...
testshort:
	go test -short -race -timeout 1m -covermode=atomic -coverprofile cover.out ./...
coverreport:
	go tool cover -html=cover.out -o cover.html
	open cover.html
clean:
	go clean
	rm -f $(BINARY_NAME)
