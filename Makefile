include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

SHELL := /bin/bash
PKG := github.com/Clever/clever-cli
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /bin | grep -v /deb | grep -v /build)
VERSION := $(shell cat VERSION)
EXECUTABLE := clever
BUILDS := \
	build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-linux-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-windows-amd64
COMPRESSED_BUILDS := $(BUILDS:%=%.tar.gz)
RELEASE_ARTIFACTS := $(COMPRESSED_BUILDS:build/%=release/%)

.PHONY: test $(PKGS) clean run vendor

$(eval $(call golang-version-check,1.9))

test: $(PKGS)

$(PKGS): golang-test-all-deps
	$(call golang-test-all,$@)


build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o "$@/$(EXECUTABLE)"
build/$(EXECUTABLE)-v$(VERSION)-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o "$@/$(EXECUTABLE)"
build/$(EXECUTABLE)-v$(VERSION)-windows-amd64:
	GOARCH=amd64 GOOS=windows go build -o "$@/$(EXECUTABLE).exe"
build: $(BUILDS)
%.tar.gz: %
	tar -C `dirname $<` -zcvf "$<.tar.gz" `basename $<`
$(RELEASE_ARTIFACTS): release/% : build/%
	mkdir -p release
	cp $< $@
release: $(RELEASE_ARTIFACTS)

clean:
	rm -rf build release

run:
	@go run main.go


install_deps: golang-dep-vendor-deps
	$(call golang-dep-vendor)