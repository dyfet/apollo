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

	"github.com/gofiber/fiber/v2"

	ipc "apollo/internal"
	"gitlab.com/tychosoft/service"
)

func addLine(ctx *fiber.Ctx) error {
	id, line := ipc.NewLine()
	if line == nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("No lines available")
	}

	lock.RLock()
	defer lock.RUnlock()
	err := ctx.Render("add-line", fiber.Map{
		"page": config,
		"Id":   id,
		"Line": line,
	})
	if err != nil {
		service.Error(err)
	}
	return err
}

func editLine(ctx *fiber.Ctx) error {
	tag := ctx.Params("id")
	if tag == "new" {
		return addLine(ctx)
	}

	id, err := strconv.Atoi(tag)
	if err != nil {
		id = 0
	}

	line := ipc.GetLine(id)
	if line == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("Line is invalid")
	}

	form := "edit-line"
	if !line.Editable {
		form = "show-line"
	}

	lock.RLock()
	defer lock.RUnlock()
	err = ctx.Render(form, fiber.Map{
		"page": config,
		"Id":   id,
		"Line": line,
	})
	if err != nil {
		service.Error(err)
	}
	return err
}

func editSettings(ctx *fiber.Ctx) error {
	lock.RLock()
	defer lock.RUnlock()
	err := ctx.Render("settings", fiber.Map{
		"page": config,
	})
	if err != nil {
		service.Error(err)
	}
	return err
}
