# Copyright (C) 2022 Tycho Softworks.
#
# This file is free software; as a special exception the author gives
# unlimited permission to copy and/or distribute it, with or without
# modifications, as long as this notice is preserved.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY, to the extent permitted by law; without even the
# implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

.PHONY: lint vet fix test cover release

ifndef	BUILD_MODE
BUILD_MODE := default
endif

GOVER=$(shell grep ^go <go.mod)
TARGET := $(CURDIR)/target
export GOCACHE := $(TARGET)/cache
export PATH := $(TARGET)/debug:${PATH}

docs:	required
	@rm -rf target/docs
	@mkdir -p target/docs
	@doc2go -out target/docs ./...

lint:	required
	@go fmt ./...
	@go mod tidy
	@staticcheck ./...

vet:	required
	@go vet ./...
	@govulncheck ./...

fix:	required
	@go fix ./...

test:
	@go test ./...

cover:	vet
	@go test -coverprofile=coverage.out ./...

go.sum:	go.mod
	@go mod tidy

# if no vendor directory (clean) or old in git checkouts
vendor:	go.sum
	@if test -d .git ; then \
		rm -rf vendor ;\
		go mod vendor ;\
	elif test ! -d vendor ; then \
		go mod vendor ;\
	else \
		touch vendor ;\
	fi
