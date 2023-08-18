SHELL=/bin/bash

# Variables
BUILD_DIR ?= ./bin
CGO_ENABLED ?= 0
VERSION ?= $(shell awk -F'"' 'NR==2{print $$4}' ./release.json)
GOOS ?= unset
GOARCH ?= unset
## If GOOS=windows, then the executable is named cybr.exe
## and the compressed file is named cybr-cli_${GOOS}_${GOARCH}.zip
ifeq (${GOOS}, windows)
TARGET_EXEC := cybr.exe
FILE_COMPRESSED := ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip
FILE_MD5 := ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip.md5
FILE_SHA256 := ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip.sha256
endif
TARGET_EXEC ?= cybr
FILE_MD5 ?= ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.md5
FILE_SHA256 ?= ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.sha256
FILE_COMPRESSED ?= ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.tar.gz


# Rules
## Version
version:
	@echo "Current version: ${VERSION}"
	@echo "GOOS: ${GOOS}"
	@echo "GOARCH: ${GOARCH}"

## Build
build:
	@echo "Building ${TARGET_EXEC} into ${BUILD_DIR}..."
	@go build -o ${BUILD_DIR}/${TARGET_EXEC} .

## Test All
test-all: vet test

## Vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## Test
test:
	@echo "Running tests..."
	@go test -v ./...

## Compile
compile:
### If GOOS is not defined, then throw an error
ifeq (${GOOS}, unset)
	@echo "GOOS is undefined. If you ran 'make release', run './release.sh' instead."
	@exit 1
endif
### If GOARCH is not defined, then throw an error
ifeq (${GOARCH}, unset)
	@echo "GOARCH is undefined"
	@exit 1
endif
	@echo "Building ${TARGET_EXEC} into ${BUILD_DIR}/${GOOS}/${GOARCH}..."
	@mkdir -p ${BUILD_DIR}/${GOOS}/${GOARCH}
	@CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BUILD_DIR}/${GOOS}/${GOARCH}/${TARGET_EXEC} .

## Install
install:
	@echo "Installing ${TARGET_EXEC} into ${HOME}/bin..."
	@mkdir -p ${HOME}/bin
	@go build -o ${HOME}/bin/${TARGET_EXEC} .

## Release
release: compile
### If GOOS=darwin, then sign the executable
ifeq (${GOOS}, darwin)
	@echo "Signing ${TARGET_EXEC} for ${GOOS} ${GOARCH}..."
	@codesign --deep --force --options=runtime \
		--sign "81176FD6CB590056A09A92B6FA9502A75F0BB3A1" \
		--timestamp "${BUILD_DIR}/${GOOS}/${GOARCH}/${TARGET_EXEC}"
	@mkdir -p ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli
endif
### If GOOS=windows, then compress the executable into a zip
ifeq (${GOOS}, windows)
	@echo "Compressing ${TARGET_EXEC} into ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip..."
	@zip -j ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip ${BUILD_DIR}/${GOOS}/${GOARCH}/${TARGET_EXEC}
	@echo "Generating MD5 checksum for ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip..."
	@md5 -qs ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip > ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip.md5
	@echo "Generating SHA256 checksum for ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip..."
	@shasum -a 256 ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip | cut -f1 -d' ' > ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.zip.sha256
### If GOOS is not windows, then compress the executable into a tar.gz
else
	@echo "Compressing ${TARGET_EXEC} into ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.tar.gz..."
	@tar -czf ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.tar.gz -C ${BUILD_DIR}/${GOOS}/${GOARCH} ${TARGET_EXEC}
	@echo "Generating MD5 checksum for ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.tar.gz..."
	@md5 -qs ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.tar.gz > ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.tar.gz.md5
	@echo "Generating SHA256 checksum for ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.tar.gz..."
	@shasum -a 256 ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.tar.gz | cut -f1 -d' ' > ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.tar.gz.sha256
endif
### If GOOS=darwin, then build, sign, and notarize a package installer
ifeq (${GOOS}, darwin)
	@echo "Building package installer structure..."
	@ditto "${BUILD_DIR}/${GOOS}/${GOARCH}/${TARGET_EXEC}" "${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli/usr/local/bin/${TARGET_EXEC}"
	@echo "Building package installer..."
	@pkgbuild --root "${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli" \
		--identifier "com.github.infamousjoeg.cybr-cli" \
		--version "${VERSION}" \
		--sign "42510BD69E802006418644C21E229DC4461D6673" \
		${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.pkg
	@echo "Notarizing package installer..."
	@xcrun notarytool submit ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.pkg \
		--keychain-profile "notary-cybrcli" \
		--wait
	@echo "Stapling package installer..."
	@xcrun stapler staple ${BUILD_DIR}/${GOOS}/${GOARCH}/cybr-cli_${GOOS}_${GOARCH}.pkg
endif

## Generate Documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	@go run docs/main.go

## Clean
clean: clean-build clean-install

clean-all: clean-build clean-install clean-compile

clean-build:
	@echo "Cleaning ${BUILD_DIR}..."
	@rm -f ${BUILD_DIR}/${TARGET_EXEC}

clean-install:
	@echo "Cleaning ${HOME}/bin..."
	@rm -f ${HOME}/bin/${TARGET_EXEC}

clean-compile:
	@echo "Cleaning ${BUILD_DIR}/linux..."
	@rm -rf ${BUILD_DIR}/linux/*
	@echo "Cleaning ${BUILD_DIR}/darwin..."
	@rm -rf ${BUILD_DIR}/darwin/*
	@echo "Cleaning ${BUILD_DIR}/windows..."
	@rm -rf ${BUILD_DIR}/windows/*