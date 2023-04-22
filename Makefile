GOCMD=go
GOBUILD=$(GOCMD) build
GOENV=$(GOCMD) env
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=sptfy
BUILD_PATH=build
BINARY_UNIX=$(BINARY_NAME)_unix
BUILD_DIR=./bin/main

all: clean build run
build:
	$(GOBUILD) -o $(BUILD_PATH)/$(BINARY_NAME) -v $(BUILD_DIR)
#test:
#	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_PATH)
run:
	./$(BUILD_PATH)/$(BINARY_NAME)
