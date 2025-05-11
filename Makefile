# Copyright (C) 2023 Tycho Softworks.
#
# This file is free software; as a special exception the author gives
# unlimited permission to copy and/or distribute it, with or without
# modifications, as long as this notice is preserved.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY, to the extent permitted by law; without even the
# implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

# Project constants
PROJECT := apollo
VERSION := 0.3.0
PATH	:= $(PWD)/target/test:${PATH}
TESTDIR := $(PWD)/web

# Debug build detects
DETECT_COVENTRY = $(shell .make/varlib.sh $(PWD)/web coventry)
DETECT_BORDEAUX = $(shell .make/varlib.sh $(PWD)/web bordeaux)

.PHONY: all required build build-test debug release install clean verify

all:            build           # default target debug
required:       vendor          # required to build
verify:		test		# verify builds
build:		lint build-test	# build defaults

# Define or override custom env
sinclude custom.mk

build-test:	required
	@install -d target/test
	@CGO_ENABLED=1 $(GO) build -v -tags debug,$(TAGS) -ldflags '-X main.etcPrefix=$(TEST_CONFIG) -X main.mediaData=$(DETECT_BORDEAUX) -X main.workingDir=$(DETECT_COVENTRY) -X main.appDataDir=$(TEST_APPDIR) -X main.logPrefix=$(TEST_LOGDIR)' -mod vendor -o target/test ./...

debug:	required
	@install -d target/debig
	@CGO_ENABLED=1 $(GO) build -v -mod vendor -tags debug,$(TAGS) -ldflags '-X main.mediaData=$(LOCALSTATEDIR)/lib/bordeaux -X main.etcPrefix=$(SYSCONFDIR) -X main.workingDir=$(LOCALSTATEDIR)/lib/coventry -X main.appDataDir=$(APPDATADIR) -X main.logPrefix=$(LOGPREFIXDIR)' -o target/debug ./...

release:	required
	@install -d target/release
	@CGO_ENABLED=1 $(GO) build --buildmode=$(BUILD_MODE) -v -mod vendor -tags release,$(TAGS) -ldflags '-s -w -X main.mediaData=$(LOCALSTATEDIR)/lib/bordeaux -X main.etcPrefix=$(SYSCONFDIR) -X main.workingDir=$(LOCALSTATEDIR)/lib/coventry -X main.appDataDir=$(APPDATADIR) -X main.logPrefix=$(LOGPREFIXDIR)' -o target/release ./...

# We normally install commandit to a local ~/go/bin for portable tooling
install:        release
	@install -d -m 755 $(DESTDIR)$(WORKINGDIR)
	@install -d -m 755 $(DESTDIR)$(SYSCONFDIR)
	@install -d -m 755 $(DESTDIR)$(SBINDIR)
	@install -d -m 755 $(DESTDIR)$(MANDIR)/man8
	@install -s -m 755 target/release/$(PROJECT) $(DESTDIR)$(SBINDIR)
	@install -m 644 etc/$(PROJECT).conf $(DESTDIR)$(SYSCONFDIR)
	@install -m 644 etc/$(PROJECT).8 $(DESTDIR)$(MANDIR)/man8
	@install -d -m 755 $(DESTDIR)$(APPDATADIR)/assets
	@install -d -m 755 $(DESTDIR)$(APPDATADIR)/views
	@find web/assets -type f -exec install -m 644 "{}" \
		$(DESTDIR)$(APPDATADIR)/assets \;
	@find web/views -type f -exec install -m 644 "{}" \
		$(DESTDIR)$(APPDATADIR)/views \;

clean:
	@$(GO) clean ./...
	@rm -rf target *.out
	@rm -f $(PROJECT)-*.tar.gz $(PROJECT)-*.tar
	@rm -f *.log

# Optional make components we add
sinclude .make/*.mk

