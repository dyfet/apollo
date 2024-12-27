#!/bin/sh
# Copyright (C) 2024 Tycho Softworks.
#
# This file is free software; as a special exception the author gives
# unlimited permission to copy and/or distribute it, with or without
# modifications, as long as this notice is preserved.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY, to the extent permitted by law; without even the
# implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

repo="$1"
verify=`git remote -v | grep push | head -1 | grep ${1}`
if test -z "$verify" ; then
	echo "*** not in release repo" ; exit 1 ; fi

status=`git status | grep "^Your branch is up to date"`
if test -z "$status" ; then
    echo "*** not in release state"
    exit 2
fi

current=`git branch -a | grep "^[*] " | cut -f2 -d\   `
case "$current" in
main|master)
    ;;
*)
    echo "*** not in release branch"
    exit 3
    ;;
esac
