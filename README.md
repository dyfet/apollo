# What is Apollo?

Apollo offers a web administrative console and web integrated telephony
services for stand-alone Coventry phone systems. This is not required for
minimal Coventry use, as Apollo simply manipulates Coventry config files, but
it does makes it much easier to run a stand-alone local Coventry service if
the config is not generated or downloaded for externally managed Coventry
devices. It is also the point of non-sip based network api contact, and can
provide additional client support services such as rosters and dialing
directories for local Coventry specific endpoints like Partisipate. Apollo
requires and supplements an install of Coventry to be at all useful.

## Dependencies

Apollo is a Go application that requires Go 1.19 or later, and GNU Make to
build. Apollo interacts with Coventry thru IPC services and manipulation of
config files, so it must be co-installed on a server running Coventry to be
used. Apollo can only be used on platforms that Coventry supports, which may
include most BSD and Linux kernel based posix systems.

While most Coventry features can be manipulated over a Web ui thru Apollo, a
Coventry "custom.conf" can override these settings and produce read-only
entries in the ui. This allows for pre-set configs for things like voice mail
extensions or door phones, while still allowing the user to modify other
extensions freely.

## Distributions

Distributions of this package are provided as detached source tarballs made
from a tagged release from our public git repository or by building the dist
target. These stand-alone detached tarballs can be used to make packages for
many GNU/Linux systems, and for BSD ports. These tagged releases already
contain all vendoring. They may be used to build and install the software
directly on a target platform without internet connections.

## Installation

From a detached tarball with embedded vendor builds, make "install" is
sufficient to install the Apollo integration server on a generic posix system.
This installs to /usr/local by default, and can be overridden with a PREFIX
setting, such as ''make PREFIX=/usr install''. The Makefile also makes it
easy to cross-compile, as well as managing separate debug and release builds.
It also should be easy to integrate detached tarballs with traditional OS
packaging.

## Participation

This project is offered as free (as in freedom) software for public use and has
a public project page at https://www.gitlab.com/tychosoft/apollo which has an
issue tracker where people can submit public bug reports and a public git
repository. Patches and merge requests may be submitted in the issue tracker or
thru email. Support requests and other kinds of inquiries may also be sent thru
the tychosoft gitlab help desktop service. Other details about participation
may be found in CONTRIBUTING.md.

