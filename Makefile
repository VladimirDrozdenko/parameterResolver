# include build/private/bgo_exports.makefile
# include ${BGO_MAKEFILE}

BGO_SPACE := .
BUILDFILE_PATH := ./build/private/bgo_exports.makefile

COPY_DIR := cp -a
REDUNDANT_FILE_LOCATION := /src/ParameterResolver
RM_RF := rm -rf
RM_F := rm -f
GO_BUILD := go build -i
BRAZIL_BUILD := false

build:: test-build 			#build-linux build-freebsd build-darwin build-windows build-linux-386 build-darwin-386 build-windows-386

clean:: remove-unneeded-files

release:: clean build quick-test

# build
.PHONY: test-build
test-build:
	mkdir -p $(BGO_SPACE)$(REDUNDANT_FILE_LOCATION)
	$(COPY_DIR) $(BGO_SPACE)/common $(BGO_SPACE)$(REDUNDANT_FILE_LOCATION)
	$(COPY_DIR) $(BGO_SPACE)/ssm $(BGO_SPACE)$(REDUNDANT_FILE_LOCATION)
	$(COPY_DIR) $(BGO_SPACE)/resolver $(BGO_SPACE)$(REDUNDANT_FILE_LOCATION)
	@echo "Building exe files..."
	go build $(BGO_SPACE)/api_driver.go
	go build $(BGO_SPACE)/cmdtool.go

# clean
.PHONY: remove-unneeded-files
remove-unneeded-files:
	$(RM_RF) build/*
	$(RM_RF) bin/
	$(RM_RF) pkg/
	$(RM_RF) vendor/bin/
	$(RM_RF) vendor/pkg/
	$(RM_RF) .cover/
	$(RM_RF) logging-files/
	$(RM_RF) src/*
	$(RM_F) api_driver
	$(RM_F) cmdtool
	find . -type f -name '*.log' -delete
	rm -rf $(BGO_SPACE)/bin/prepacked

# release
.PHONY: quick-test
quick-test:
	go test -gcflags "-N -l" -tags=integration parameterResolver/resolver
	go test -gcflags "-N -l" -tags=integration parameterResolver/ssm
