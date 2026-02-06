# Detect OS and set binary name
OS := $(shell uname -s)
ifeq ($(findstring Darwin, $(OS)), Darwin)
	BINARY_NAME = hzmind
else ifeq ($(findstring Linux, $(OS)), Linux)
	BINARY_NAME = hzmind
else
	BINARY_NAME = hzmind.exe
endif

BUILD_DATE = $(shell date +%Y%m%d%H%M)

MAIN_PATH = cmd/hzmind/main.go

BUILD_DIR = ./bin
INTERNAL_DIR = ./internal
INSTALL_DIR = $(APPDATA)/HarzMindCode

VERSION_FILE = version.go

.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/n] ' && read ans && [ $${ans:-n} = y ]

.PHONY: genver
genver:
	@rm -f ${INTERNAL_DIR}/$(VERSION_FILE)
	@echo "package internal" > ${INTERNAL_DIR}/$(VERSION_FILE)
	@echo "" >> ${INTERNAL_DIR}/$(VERSION_FILE)
	@echo "const VERSION_DATE = \"$(BUILD_DATE)\"" >> ${INTERNAL_DIR}/$(VERSION_FILE)

## build: build the application
.PHONY: build
build: genver
	go build -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_PATH}

## export: export the application
.PHONY: export
export: genver
	@mkdir -p ${BUILD_DIR}/${BUILD_DATE}
	#GOARCH=amd64 GOOS=darwin go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}-darwin-amd64 ${MAIN_PATH}
	#GOARCH=arm64 GOOS=darwin go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}-darwin-arm64 ${MAIN_PATH}
	#GOARCH=amd64 GOOS=linux go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}-linux-amd64 ${MAIN_PATH}
	GOARCH=amd64 GOOS=windows go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}-windows-amd64.exe ${MAIN_PATH}
	#GOARCH=amd64 GOOS=windows go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}.exe ${MAIN_PATH}

## run: run the application
.PHONY: run
run:
	${BUILD_DIR}/${BINARY_NAME} -l

## clean: clean up the build binaries
.PHONY: clean
clean: confirm
	@echo "Cleaning up..."
	@rm -rf ${BUILD_DIR}

## install: install the application for the user
.PHONY: install
install:
	@mkdir -p $(INSTALL_DIR)
	@echo "Copying binary to $(INSTALL_DIR)..."
	@cp ${BUILD_DIR}/${BINARY_NAME} $(INSTALL_DIR)/
	@echo "Installation complete! The binary has been copied to $(INSTALL_DIR)."
	@echo "To make it available in your PATH, add $(APPDATA)/HarzMindCode to your environment variables manually."
	@echo "On Windows, go to System Properties > Environment Variables and edit the User PATH."
	@echo "Please restart your terminal for changes to take effect."
