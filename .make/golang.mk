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

ifndef	GO
GO := go
endif

STATIC_CHECK	:= $(shell which staticcheck 2>/dev/null || true )
ifeq ($(STATIC_CHECK),)
STATIC_CHECK	:= true
endif

GOVULN_CHECK	:= $(shell which govulncheck 2>/dev/null || true)
ifeq ($(GOVULN_CHECK),)
GOVULN_CHECK	:= true
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
	@$(GO) fmt ./...
	@$(GO) mod tidy
	@$(STATIC_CHECK) ./...

vet:	required
	@$(GO) vet ./...
	@$(GOVULN_CHECK) ./...

fix:	required
	@$(GO) fix ./...

test:
	@$(GO) test ./...

stage:
	rm -rf target/stage
	mkdir -p target/stage
	$(MAKE) DESTDIR=$(CURDIR)/target/stage install

cover:	vet
	@$(GO) test -coverprofile=coverage.out ./...

go.sum:	go.mod
	@$(GO) mod tidy

# if no vendor directory (clean) or old in git checkouts
vendor:	go.sum
	@if test -d .git ; then \
		rm -rf vendor ;\
		$(GO) mod vendor ;\
	elif test ! -d vendor ; then \
		$(GO) mod vendor ;\
	else \
		touch vendor ;\
	fi
