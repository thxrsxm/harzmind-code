BINARY_NAME = hzmind
BUILD_DATE = $(shell date +%Y%m%d%H%M)

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
	go build -o ${BINARY_NAME}.exe main.go

## export: export the application
.PHONY: export
export: genver
	@mkdir -p ${BUILD_DIR}/${BUILD_DATE}
	#GOARCH=amd64 GOOS=darwin go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}-darwin-amd64 main.go
	#GOARCH=arm64 GOOS=darwin go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}-darwin-arm64 main.go
	#GOARCH=amd64 GOOS=linux go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}-linux-amd64 main.go
	GOARCH=amd64 GOOS=windows go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}-windows-amd64.exe main.go
	#GOARCH=amd64 GOOS=windows go build -o ${BUILD_DIR}/${BUILD_DATE}/${BINARY_NAME}.exe main.go

## run: run the application
.PHONY: run
run:
	./${BINARY_NAME}.exe

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
	@cp ${BINARY_NAME}.exe $(INSTALL_DIR)/
	@echo "Installation complete! The binary has been copied to $(INSTALL_DIR)."
	@echo "To make it available in your PATH, add $(APPDATA)\HarzMindCode to your environment variables manually."
	@echo "On Windows, go to System Properties > Environment Variables and edit the User PATH."
	@echo "Please restart your terminal for changes to take effect."
