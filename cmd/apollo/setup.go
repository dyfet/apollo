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
	"net"
	"strings"
	"syscall"

	"github.com/glendc/go-external-ip"
	"github.com/gofiber/fiber/v2"
	"github.com/ipinfo/go/v2/ipinfo"

	"apollo/internal"
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
	apollo.SetConfig(section, "webadmin", admin)

	iniCoventry := workingDir + "/dynamic.conf"
	err := dynCoventry.SaveTo(iniCoventry)
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

func locationSetup(ctx *fiber.Ctx) error {
	lock.Lock()
	defer lock.Unlock()
	apollo.UpdateCoventry("server", "location", ctx.FormValue("geolocated"))
	apollo.UpdateCoventry("server", "where", ctx.FormValue("where"))
	apollo.UpdateCoventry("server", "city", ctx.FormValue("city"))
	apollo.UpdateCoventry("server", "region", ctx.FormValue("region"))
	apollo.UpdateCoventry("server", "postal", ctx.FormValue("postal"))
	apollo.SaveCoventry()
	service.Debug(3, "set where ", ctx.FormValue("where"))

	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	return ctx.Redirect("/lines", fiber.StatusSeeOther)
}

func internetSetup(ctx *fiber.Ctx) error {
	pubip := ctx.FormValue("publicip")
	token := ctx.FormValue("iptoken")

	service.Info("looking up address...")
	// If offline, try finding address
	if pubip == "offline" || pubip == "down" || pubip == "auto" {
		consensus := externalip.DefaultConsensus(nil, nil)
		ip, err := consensus.ExternalIP()
		if err == nil {
			pubip = ip.String()
			service.Info("public address ", pubip)
		}
	}

	// Process token

	if pubip == "offline" || pubip == "down" || pubip == "auto" {
		publicIp = "down"
		return ctx.Redirect("/settings", fiber.StatusSeeOther)
	}
	if pubip == "offline" || pubip == "down" {
		pubip = "auto"
	}
	publicIp = pubip
	service.Info("looking up location...")

	lock.Lock()
	defer lock.Unlock()
	apollo.UpdateCoventry("server", "token", token)
	apollo.SaveCoventry()
	service.Debug(3, "set token ", token)

	client := ipinfo.NewClient(nil, nil, token)
	info, err := client.GetIPInfo(net.ParseIP(pubip))
	if err != nil {
		service.Error(err)
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		return ctx.Redirect("/lines", fiber.StatusSeeOther)

	}

	location := info.Location
	country := strings.ToLower(info.Country)
	city := info.City
	region := info.Region
	postal := info.Postal
	timezone := info.Timezone
	service.Info("located: ", location, " ", city, ", ", region, ",", postal, ", ", country)

	apollo.UpdateCoventry("server", "location", location)
	apollo.UpdateCoventry("server", "city", city)
	apollo.UpdateCoventry("server", "region", region)
	apollo.UpdateCoventry("server", "postal", postal)
	apollo.UpdateCoventry("server", "country", country)
	apollo.UpdateCoventry("weather", "timezone", timezone)
	apollo.SaveCoventry()

	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	return ctx.Redirect("/lines", fiber.StatusSeeOther)
}

func themeSetup(ctx *fiber.Ctx) error {
	server := apollo.GetServer()
	theme := apollo.GetConfig(server, "theme", "dark")
	if theme == "light" {
		theme = "dark"
	} else {
		theme = "light"
	}

	service.Debug(3, "set theme ", theme)
	apollo.UpdateCoventry("server", "theme", theme)
	apollo.SaveCoventry()
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	return ctx.Redirect("/lines", fiber.StatusSeeOther)
}
