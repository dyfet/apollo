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
	"sort"

	"github.com/gofiber/fiber/v2"

	"apollo/internal"
	"gitlab.com/tychosoft/service"
)

func viewPing(ctx *fiber.Ctx) error {
	// render master site ping page
	lock.RLock()
	defer lock.RUnlock()
	err := ctx.Render("ping", fiber.Map{
		"page": config,
	})
	if err != nil {
		service.Error(err)
	}
	return err
}

func viewMain(ctx *fiber.Ctx) error {
	// will probably have some info pulled from status
	// will have button for change password, to add line, add group, etc...
	lock.RLock()
	defer lock.RUnlock()
	err := ctx.Render("main", fiber.Map{
		"page": config,
	})
	if err != nil {
		service.Error(err)
	}
	return err
}

func viewLines(ctx *fiber.Ctx) error {
	type Item struct {
		Id   int
		Line *apollo.Line
	}

	lines := apollo.GetLines()
	items := make([]Item, 0, len(lines))
	for key, value := range lines {
		items = append(items, Item{key, value})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Id < items[j].Id
	})

	lock.RLock()
	defer lock.RUnlock()
	err := ctx.Render("lines", fiber.Map{
		"page":  config,
		"items": items,
	})
	if err != nil {
		service.Error(err)
	}
	return err
}

func viewGroups(ctx *fiber.Ctx) error {
	// will probably pass group config, have manipulation functions
	lock.RLock()
	defer lock.RUnlock()
	err := ctx.Render("groups", fiber.Map{
		"page": config,
	})
	if err != nil {
		service.Error(err)
	}
	return err
}

func viewContacts(ctx *fiber.Ctx) error {
	// will form list from config, have print option
	lock.RLock()
	defer lock.RUnlock()
	err := ctx.Render("contacts", fiber.Map{
		"page": config,
	})
	if err != nil {
		service.Error(err)
	}
	return err
}

func viewSetup(ctx *fiber.Ctx) error {
	if setupFlag {
		return ctx.Redirect("/lines", fiber.StatusTemporaryRedirect)
	}

	lock.RLock()
	defer lock.RUnlock()
	err := ctx.Render("setup", fiber.Map{
		"page":    config,
		"weather": weather,
	})
	if err != nil {
		service.Error(err)
	}
	return err
}
