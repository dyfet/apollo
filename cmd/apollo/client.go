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

	"apollo/internal"
)

func clientPing(ctx *fiber.Ctx) error {
	id := ctx.Locals("userID").(int)
	return ctx.SendString("User: " + strconv.Itoa(id))
}

func clientProfile(ctx *fiber.Ctx) error {
	id := ctx.Locals("userID").(int)
	profile := apollo.GetLine(id)
	if profile != nil {
		return ctx.JSON(profile)
	}
	return fiber.NewError(fiber.StatusNotFound, "Profile not found")
}

func clientRoster(ctx *fiber.Ctx) error {
	return ctx.JSON(apollo.GetLines())
}

func clientGroups(ctx *fiber.Ctx) error {
	return ctx.JSON(apollo.GetGroups())
}
