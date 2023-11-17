GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=myriad
BUILD_DIR=./cmd/myriad
MYRIAD_SAMPLE_FILE=./test/sample.my

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(BUILD_DIR)
clean:
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME) $(MYRIAD_SAMPLE_FILE)