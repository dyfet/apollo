# Features

Apollo provides a web ui to manage and manipulate a co-resident Coventry
server. This web ui makes minimal use of Javascript, and may be navigable on
non-JS browsers as well. The ui is primarily defined in web templates and forms
which can easily be re-written for multi-language support. Light and dark
themes are available and the ui uses high contrast colors to help make it
easier for visually impaired users. Borders are often outlined for similar
reasons. The ui goals are simplicity and discoverability.

## 1. Device Setup

Apollo interacts with Coventry (and Bordeaux) thru IPC services (posix message
queues and shared memory) and by manipulating config files. It can provide web
api services for Apollo as well as ui access. It also maintains a web "user" id
for login acess. When initially installed, the web ui offers a setup page,
which can be used to initialize the web user id and test ipc operations.

Other global service attributes can be specified as part of this setup,
including the theme to be used, and the sip admininstrative domain that
coventry will operate under. The sip administrative domain cannot be changed
after setup as it's part of the password hash of each sip extension.

Part of how Apollo configures Coventry is thru the Coventry "dynamic.conf".
This file is created by Apollo when it is started. It is used to hold all
Apollo config "changes" such as global configs, extension lines, passwords, and
the web admin user. You can restore Apollo to it's initial install state simply
by removing it.

## 2. Line Management

Once setup you always login to the line management screen. This shows you what
current extension lines you have defined and a common navigation bar at the top
used in all Apollo main screens. This lets you see their presense and
registration state, display name, line #, and similar properties.  A more
detailed view of an individual extension can be shown. This also lets you add
an extension, modify an existing entry such as to change it's display name or
registration password, or to remote it.

## 3. Group Management

Groups are a set of lines that can be accessed thru a virtual (3 digit or
larger) extension number. A single line may be a laptop, and another may be a
desktop phone, but by making a group, you can form a single virtual extension
number that can then ring both. Some may prefer to do internal dial plans
entirely under group management, and use device lines simply for direct device
ringing and internal management.

Group behavior will also allow for scheduling of immediate and delayed ring
call coverage. This behavior is part of Coventry's config already, but is not
yet fully exposed in the web ui. So, for example, a primary user's phone may
ring immediately, and coverage positions can start ringing later if it is not
picked up. Voice mail can then be implimented as a later ringing call coverage
operation in Coventry, rather than thru traditional call forwarding.

## 4. Contact Management

Contact management includes personal (per user) and global speed dialing. It
will also eventually include directories for external numbers.

## 5. Settings Management

These are global settings that can be changed without reseting all Apollo's
dynamic.conf entries. This typically includes the current theme and the web
admin password. It can in the future include info on upstream dialing
providers.

## 6. Trouble Ticketing

Trouble ticketing and basic tracking will be added. This will include basic
cable plan management as well.

## 7. Call Accounting

Call accounting of external traffic is to be added.

## 8. User Page

A user "page" may become available where a user enters their line extension and
password, and can the access personal settings they can directly maintain on
their own, and issue service requests.

## 9. Api services

Pure web api services will be added thru Apollo.
