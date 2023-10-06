ifeq ($(OS),Windows_NT)
    BINARY_SUFFIX=.exe
else
    BINARY_SUFFIX=
endif

build:
	@go build -o bin/git-auto-config$(BINARY_SUFFIX) auto-config.go
	@go build -o bin/git-clip$(BINARY_SUFFIX) clip.go
	@go build -o bin/git-search$(BINARY_SUFFIX) search.go
