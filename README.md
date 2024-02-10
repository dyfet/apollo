# What is Apollo?

Apollo offers a web administrative console and web integrated telephony
services for stand-alone Coventry phone systems. This is not required for
minimal Coventry use, as Apollo simply manipulates Coventry config files, but
it does makes it much easier to run a stand-alone local Coventry service if
the config is not generated or downloaded for externally managed Coventry
devices. It is also the point of non-sip based network api contact, and can
provide additional client support services such as rosters and dialing
directories for local Coventry specific endpoints like Partisipate. Apollo
requires and suppliments an install of Coventry to be at all useful.

## Dependencies

Apollo is a Go application that requires Go 1.19 or later, and GNU Make to
build. Apollo interacts with Coventry thru IPC services and manipulation of
config files, so it must be co-installed on a server running Coventry to be
used. Apollo can only be used on platforms that Coventry supports, which may
include NetBSD (10), FreeBSD, and most Linux kernel based distributions.

While most Coventry features can be manipulated over a Web ui thru Apollo, a
Coventry "custom.conf" can override these settings and produce read-only
entries in the ui. This allows for pre-set configs for things like voice mail
extensions or door phones, while still allowing the user to modify other
extenstions freely.

## Distributions

Distributions of this package are provided as detached source tarballs made
from a tagged release from our internal source repository. These stand-alone
detached tarballs can be used to make packages for many GNU/Linux systems, and
for BSD ports. They may also be used to build and install the software
directly on a target platform.

The source tarballs bundle a vendor directory. This has a local locked copy of
ALL third party Go imports. This allows you to take the detached source
tarball, and build an Apollo service identically to how I may do so for binary
distributions, even when doing builds in a secure offline build environment.

The latest release source tarball is found at
https://www.tychosoft.com/tychosoft/-/packages/generic/apollo which provides
access to past releases as well.

## Installation

Make "install" is sufficient to install the Apollo integration server on a
generic posix system. This installs to /usr/local by default, and can be
overridden with a PREFIX setting, such as ''make PREFIX=/usr install''. The
Makefile also makes it easy to cross-compile, as well as managing separate
debug and release builds. It also should be easy to integrate with traditional
OS packaging.

## Participation

This project is offered as free (as in freedom) software for public use and has
a public home page at https://www.tychosoft.com/tychosoft/apollo which has an
issue tracker where people can submit public bug reports, and a wiki for
hosting project documentation. We are not maintaining a public git repo nor
do we have any production or development related resources hosted on external
sites. Patches may be submitted and attached to an issue in the issue tracker.
Support requests and other kinds of inquiries may also be sent privately thru
email to tychosoft@gmail.com. Other details about participation may be found
in the Contributing page.
