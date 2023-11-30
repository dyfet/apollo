#!/bin/sh
# Copyright (C) 2023 Tycho Softworks.
#
# This file is free software; as a special exception the author gives
# unlimited permission to copy and/or distribute it, with or without
# modifications, as long as this notice is preserved.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY, to the extent permitted by law; without even the
# implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

if test -d /var/lib/$2 ; then
	echo /var/lib/$2 ; exit 0 ; fi

if test -d /usr/var/lib/$2 ; then
	echo /usr/var/lib/$2 ; exit 0 ; fi

if test -d /usr/local/var/lib/$2 ; then
	echo /usr/local/var/lib/$2 ; exit 0 ; fi

if test -d /usr/pkg/var/lib/$2 ; then
	echo /usr/pkg/var/lib/$2 ; exit 0 ; fi

if test -d ../$2/test ; then
	echo $(pwd)/../$2/test ; exit 0 ; fi

if test -d ../../build/$2/test ; then
	echo $(pwd)/../../build/$2/test ; exit 0 ; fi

echo $1
exit 0
