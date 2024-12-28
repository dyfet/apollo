// Copyright (C) 2023 Tycho Softworks.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"

	"apollo/internal"
	"gitlab.com/tychosoft/service"
)

type User struct {
	Username string
	Password string
}

var (
	dynCoventry *ini.File = nil
	adminUser   *User
	setupFlag   bool = true
)

func dynInit(port uint16, tls bool) {
	var opt = ini.LoadOptions{Loose: true, Insensitive: true}
	var err error
	var web = fmt.Sprintf("%d", port)
	var iniCoventry = workingDir + "/dynamic.conf"

	dynCoventry, err = ini.LoadSources(opt, iniCoventry)
	if err != nil {
		dynCoventry = &ini.File{}
	}
	section := dynCoventry.Section("server")
	apollo.SetConfig(section, "webserver", web)
	if tls {
		apollo.SetConfig(section, "urlschema", "https")
	} else {
		apollo.SetConfig(section, "urlschema", "http")
	}

	if !section.HasKey("webadmin") {
		section.NewKey("webadmin", "admin")
	}

	setupFlag = section.HasKey("webpass")
	adminUser = &User{
		Username: apollo.GetConfig(section, "webadmin", "admin"),
		Password: apollo.GetConfig(section, "webpass", "XXX"),
	}
	err = dynCoventry.SaveTo(iniCoventry)
	os.Chmod(iniCoventry, 0600)
	if err == nil {
		err = apollo.ReloadCoventry()
	}
	if err != nil {
		service.Error(err)
	}
}
