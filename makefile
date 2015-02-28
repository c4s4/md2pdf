NAME=md2xml
BUILD_DIR=build
TEST_DIR=test

YELLOW=\033[93m
RED=\033[91m
CLEAR=\033[0m


all: clean test

bin:
	@echo "$(YELLOW)Building executable$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go build md2xml.go
	mv md2xml $(BUILD_DIR)
	go build md2pdf.go
	mv md2pdf $(BUILD_DIR)

test: bin
	@echo "$(YELLOW)Running test$(CLEAR)"
	$(BUILD_DIR)/md2xml -o $(BUILD_DIR)/example.xml $(TEST_DIR)/example.md
	xmllint --noout $(BUILD_DIR)/example.xml
	$(BUILD_DIR)/md2pdf -o $(BUILD_DIR)/example.pdf -i $(TEST_DIR) $(TEST_DIR)/example.md

install: bin
	@echo "$(YELLOW)Installing binaries$(CLEAR)"
	sudo cp $(BUILD_DIR)/md2xml $(BUILD_DIR)/md2pdf /opt/bin/

clean:
	@echo "$(YELLOW)Cleaning generated files$(CLEAR)"
	rm -rf $(BUILD_DIR)
