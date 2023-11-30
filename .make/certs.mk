# Copyright (C) 2023 Tycho Softworks.
#
# This file is free software; as a special exception the author gives
# unlimited permission to copy and/or distribute it, with or without
# modifications, as long as this notice is preserved.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY, to the extent permitted by law; without even the
# implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

.PHONY:	certs

certs:
	@rm -f web/server.key web/server.crt
	@openssl ecparam -genkey -name secp384r1 -out web/server.key
	@openssl req -new -x509 -sha256 -key web/server.key -out web/server.crt -days 3650

