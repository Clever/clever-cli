SHELL := /bin/bash
PKG = github.com/Clever/clever-cli
SUBPKGSREL := $(shell ls -d */ | grep -v bin | grep -v deb)
SUBPKGS = $(addprefix $(PKG)/,$(SUBPKGSREL))
PKGS = $(PKG) $(SUBPKGS)

.PHONY: test golint

golint:
	@go get github.com/golang/lint/golint

test: $(PKGS)

$(PKGS): golint
	@go get -d -t $@
	@gofmt -w=true $(GOPATH)/src/$@*/**.go
ifneq ($(NOLINT),1)
	@echo "LINTING..."
	@PATH=$(PATH):$(GOPATH)/bin golint $(GOPATH)/src/$@*/**.go
	@echo ""
endif
ifeq ($(COVERAGE),1)
	@go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	@go tool cover -html=$(GOPATH)/src/$@/c.out
else
	@echo "TESTING..."
	@go test $@ -test.v
endif

run:
	@go run main.go
