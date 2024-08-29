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
	"crypto/sha256"
	"encoding/hex"
	"syscall"

	"github.com/gofiber/fiber/v2"

	ipc "apollo/internal"
	"gitlab.com/tychosoft/service"
)

func postSetup(ctx *fiber.Ctx) error {
	// Make sure we do nothing if already setup
	if setupFlag {
		return ctx.Status(fiber.StatusBadRequest).SendString("Apollo is already configured")
	}

	admin := ctx.FormValue("admin")
	passwd := ctx.FormValue("pass")
	verify := ctx.FormValue("verify")

	// Server side verification repeats javascript tests for non-js browsers
	if admin == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Please enter an admin username.")
	}

	if passwd == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Please enter a password.")
	}

	if passwd != verify {
		return ctx.Status(fiber.StatusBadRequest).SendString("Password does not match verify.")
	}

	// Process form and goto main...
	lock.Lock()
	defer lock.Unlock()

	section := dynCoventry.Section("server")
	digest := sha256.New()
	digest.Write([]byte(passwd + ":" + admin))
	passwd = hex.EncodeToString(digest.Sum(nil))
	section.NewKey("webpass", passwd)
	ipc.SetConfig(section, "webadmin", admin)

	covpath := workingDir + "/dynamic.conf"
	err := dynCoventry.SaveTo(covpath)
	if err != nil {
		service.Error(err)
		return ctx.Status(fiber.StatusBadRequest).SendString("Cannot save setup")
	}

	adminUser = &User{
		Username: admin,
		Password: passwd,
	}

	setupFlag = true
	return ctx.Redirect("/lines", fiber.StatusSeeOther)
}

func themeSetup(ctx *fiber.Ctx) error {
	server := ipc.GetServer()
	theme := ipc.GetConfig(server, "theme", "dark")
	if theme == "light" {
		theme = "dark"
	} else {
		theme = "light"
	}

	ipc.UpdateCoventry("server", "theme", theme)
	ipc.SaveCoventry()
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	return ctx.Redirect("/lines", fiber.StatusSeeOther)
}
