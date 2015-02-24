NAME=md2xml
BUILD_DIR=build
TEST_DIR=test

YELLOW=\033[93m
RED=\033[91m
CLEAR=\033[0m


all: clean test

bin:
	@echo "$(YELLOW)Building executable$(CLEAR)"
	go build md2xml.go
	mkdir -p $(BUILD_DIR)
	mv md2xml $(BUILD_DIR)

test: bin
	@echo "$(YELLOW)Running test$(CLEAR)"
	$(BUILD_DIR)/md2xml -o $(BUILD_DIR)/example.xml $(TEST_DIR)/example.md
	xmllint --noout $(BUILD_DIR)/example.xml

clean:
	@echo "$(YELLOW)Cleaning generated files$(CLEAR)"
	rm -rf $(BUILD_DIR)
