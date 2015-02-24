NAME=md2xml
BUILD_DIR=build

YELLOW=\033[93m
RED=\033[91m
CLEAR=\033[0m


all: test

binary:
	@echo "$(YELLOW)Building executable$(CLEAR)"
	go build md2xml.go
	mkdir -p $(BUILD_DIR)
	mv md2xml $(BUILD_DIR)

test: binary
	@echo "$(YELLOW)Running test$(CLEAR)"
	$(BUILD_DIR)/md2xml -o $(BUILD_DIR)/example.xml test/example.md
	xmllint --noout $(BUILD_DIR)/example.xml
