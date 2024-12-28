# Copyright (C) 2022 Tycho Softworks.
#
# This file is free software; as a special exception the author gives
# unlimited permission to copy and/or distribute it, with or without
# modifications, as long as this notice is preserved.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY, to the extent permitted by law; without even the
# implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

.PHONY: dist distvet distclean mindist

dist:	required
	@rm -f $(PROJECT)-*.tar.gz $(PROJECT)-*.tar
	@git archive -o $(PROJECT)-$(VERSION).tar --format tar --prefix=$(PROJECT)-$(VERSION)/ v$(VERSION) 2>/dev/null || git archive -o $(PROJECT)-$(VERSION).tar --format tar --prefix=$(PROJECT)-$(VERSION)/ HEAD
	@if test -f vendor/modules.txt ; then \
		tar --transform s:^:$(PROJECT)-$(VERSION)/: --append --file=$(PROJECT)-$(VERSION).tar vendor ; fi
	@gzip $(PROJECT)-$(VERSION).tar

mindist:
	@rm -f $(ARCHIVE)-*.tar.gz $(ARCHIVE)-*.tar
	@git archive -o $(ARCHIVE)-mindist.tar --format tar --prefix=$(ARCHIVE)-mindist/ HEAD
	@gzip $(ARCHIVE)-mindist.tar

distvet:	vet

distclean:	clean
	@if test -f go.mod ; then rm -rf vendor ; fi
	@rm -f go.sum
	@if test -f go.mod ; then ( echo "module $(PROJECT)" ; echo ; echo "$(GOVER)" ) > go.mod ; fi
