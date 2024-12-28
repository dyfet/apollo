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
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"apollo/internal"
)

func deleteLine(ctx *fiber.Ctx) error {
	ext := ctx.Params("id")
	id, err := strconv.Atoi(ext)
	if err != nil {
		id = 0
	}

	if ext != ctx.FormValue("line") {
		return ctx.Status(fiber.StatusBadRequest).SendString("line does not match id")
	}

	apollo.RemoveLine(id)
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	return ctx.Redirect("/lines", fiber.StatusSeeOther)
}

func passwdLine(ctx *fiber.Ctx) error {
	ext := ctx.Params("id")
	id, err := strconv.Atoi(ext)
	if err != nil {
		id = 0
	}

	line := apollo.GetLine(id)
	if line == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("Line is invalid")
	}

	if !line.Editable {
		return ctx.Status(fiber.StatusBadRequest).SendString("Custom lines not changeable")
	}

	save := apollo.SavedLine(id)
	passwd := ctx.FormValue("pass")
	verify := ctx.FormValue("verify")

	if passwd != verify {
		return ctx.Status(fiber.StatusBadRequest).SendString("Password does not match verify.")
	}

	if len(passwd) == 0 {
		return ctx.Status(fiber.StatusBadRequest).SendString("Password not set.")
	}

	save.MD5 = ""
	save.SHA256 = ""
	save.Secret = ""

	if apollo.HasMD5() {
		save.MD5 = apollo.ComputeMD5(ext, passwd)
	}

	if apollo.HasSHA256() {
		save.SHA256 = apollo.ComputeSHA256(ext, passwd)
	}

	apollo.UpdateLine(id, save)
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	return ctx.Redirect("/lines", fiber.StatusSeeOther)
}

func postNewLine(ctx *fiber.Ctx) error {
	ext := ctx.FormValue("ext")
	id, _ := strconv.Atoi(ext)

	if apollo.ExistsLine(id) {
		return ctx.Status(fiber.StatusBadRequest).SendString("Line already exists")
	}

	save := &apollo.Line{}
	save.Type = ctx.FormValue("type")
	save.Display = ctx.FormValue("display")

	count := ctx.FormValue("lines")
	lines, err := strconv.Atoi(count)
	save.Lines = uint16(lines)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	passwd := ctx.FormValue("newp")
	verify := ctx.FormValue("verify")

	if passwd != verify {
		return ctx.Status(fiber.StatusBadRequest).SendString("Password does not match verify.")
	}

	if len(passwd) > 0 {
		if apollo.HasMD5() {
			save.MD5 = apollo.ComputeMD5(ext, passwd)
		}

		if apollo.HasSHA256() {
			save.SHA256 = apollo.ComputeSHA256(ext, passwd)
		}
	}

	apollo.UpdateLine(id, save)
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	return ctx.Redirect("/lines", fiber.StatusSeeOther)
}

func postLine(ctx *fiber.Ctx) error {
	ext := ctx.Params("id")
	id, err := strconv.Atoi(ext)
	if err != nil {
		id = 0
	}

	line := apollo.GetLine(id)
	if line == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("Line is invalid")
	}

	if !line.Editable {
		return ctx.Status(fiber.StatusBadRequest).SendString("Custom lines not changeable")
	}

	save := apollo.SavedLine(id)

	if ctx.FormValue("type") != line.Type {
		save.Type = ctx.FormValue("type")
	}

	if ctx.FormValue("display") != line.Display {
		save.Display = ctx.FormValue("display")
	}

	if ctx.FormValue("caller") != line.Caller {
		save.Caller = ctx.FormValue("caller")
	}

	if ctx.FormValue("email") != line.EMail {
		save.EMail = ctx.FormValue("email")
	}

	if ctx.FormValue("cabling") != line.Cabling {
		save.Cabling = ctx.FormValue("cabling")
	}

	if ctx.FormValue("location") != line.Location {
		save.Location = ctx.FormValue("location")
	}

	count := ctx.FormValue("lines")
	lines, err := strconv.Atoi(count)
	if uint16(lines) != line.Lines {
		save.Lines = uint16(lines)
	}

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	apollo.UpdateLine(id, save)
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	return ctx.Redirect("/lines", fiber.StatusSeeOther)
}
