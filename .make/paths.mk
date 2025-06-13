# Copyright (C) 2024 Tycho Softworks.
#
# This file is free software; as a special exception the author gives
# unlimited permission to copy and/or distribute it, with or without
# modifications, as long as this notice is preserved.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY, to the extent permitted by law; without even the
# implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

PROJECT ?= $(shell basename `git rev-parse --show-toplevel`)
ifeq ($(OS),Windows_NT)
OUTPUT  := $(PROJECT).exe
else
OUTPUT  := $(PROJECT)
endif

# Project overrides, starting with prefix install
TAGS =

ifeq ($(OS),Windows_NT)
ifndef	PREFIX
PREFIX := "C:\\Program Files\\$(PROJECT)"
endif

ifndef	SYSCONFDIR
SYSCONFDIR := "C:\\ProgramData\\$(PROJECT)"
endif

ifndef	WORKINGDIR
WORKINGDIR := "C:\\ProgramData\\$(PROJECT)"
endif

ifndef	LOCALSTATEDIR
LOCALSTATEDIR := "C:\\ProgramData\\$(PROJECT)"
endif

ifndef	RUNPREFIXDIR
RUNPREFIXDIR := "C:\\ProgramData\\$(PROJECT)"
endif
endif

ifndef	DESTDIR
DESTDIR =
endif

ifndef	PREFIX
PREFIX := /usr/local
endif

ifndef	BINDIR
BINDIR := $(PREFIX)/bin
endif

ifndef	SBINDIR
SBINDIR := $(PREFIX)/sbin
endif

ifndef	LIBDIR
LIBDIR := $(PREFIX)/lib
endif

ifndef	LIBDATADIR
LIBDATADIR := $(PREFIX)/lib
endif

ifndef	DATADIR
DATADIR := $(PREFIX)/share
endif

ifndef	MANDIR
MANDIR := $(DATADIR)/man
endif

ifndef	LOCALSTATEDIR
LOCALSTATEDIR := $(PREFIX)/var
endif

ifndef	SYSCONFDIR
SYSCONFDIR := $(PREFIX)/etc
endif

ifndef	LOGPREFIXDIR
LOGPREFIXDIR := $(LOCALSTATEDIR)/log
endif

ifndef RUNPREFIXDIR
RUNPREFIXDIR := $(LOCALSTATEDOR)/run
endif

ifndef	WORKINGDIR
WORKINGDIR := $(LOCALSTATEDIR)/lib/$(PROJECT)
endif

ifndef	APPDATADIR
APPDATADIR := $(DATADIR)/$(PROJECT)
endif
