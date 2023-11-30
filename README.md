# What is Apollo?

Apollo offers a web administrative console and web integrated telephony
services for a Coventry phone system.  It is the point of network api contact,
and provides additional client support services, such as rosters and dialing
directories for coventry specific endpoints like partisipate. Apollo requires
and suppliments an install of Coventry to be at all useful.

## Installation

Make "install" is sufficient to install the Apollo integration server on a
generic posix system. This installs to /usr/local by default, and can be
overridden with a PREFIX setting, such as ''make PREFIX=/usr install''. The
Makefile also makes it easy to cross-compile, as well as managing separate
debug and release builds. It also should be easy to integrate with traditional
OS packaging.

In git checkouts I manage a vendor directory outside of git.  This is because
it may generate different content when you update the go.mod file. Generating
a vendor branch means it also can get into the stand-alone dist tarball, and
that can then be used in network isolated build systems.  Since the builds are
cached anyway without a vendor directory, this has no impact on performance.
The vendor directory is only refreshed if the go.sum file changes.

