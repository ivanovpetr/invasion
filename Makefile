BUILD_FOLDER = ./build

build:
	@echo Building Invasion...
	@-mkdir -p $(BUILD_FOLDER) 2> /dev/null
	@go build -o $(BUILD_FOLDER) ./...

lint:
	@echo Running gocilint...
	@golangci-lint run --out-format=tab --issues-exit-code=0

test:
	@echo Running tests...
	@go test -race -failfast -v -timeout 5m ./...

clean-build:
	rm -rf $(BUILD_FOLDER)
rebuild: clean-build build