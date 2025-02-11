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
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/alexflint/go-arg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/template/html/v2"
	"gopkg.in/ini.v1"

	"apollo/internal"
	"gitlab.com/tychosoft/service"
)

// Server config
type Config struct {
	// apollo config data
	Host    string `ini:"host" arg:"--host" help:"server bind"`
	Port    uint16 `ini:"port" arg:"--port" help:"server port"`
	Secure  bool   `ini:"secure" arg:"-s,--secure" help:"Server tls mode"`
	Verbose int    `ini:"verbose" help:"debugging log level (also -v..)"`

	// certificate info
	Keyfile string `ini:"keyfile" arg:"-"`
	Crtfile string `ini:"crtfile" arg:"-"`

	// page specific and branding
	Logo     string `ini:"logo" arg:"-"`
	Home     string `ini:"home" arg:"-"`
	Views    string `ini:"views" arg:"-"`
	Theme    string `ini:"theme" arg:"-"`
	Weather  string `ini:"weather" arg:"-"`
	Location string `ini:"location" arg:"-"`
	Where    string `ini:"where" arg:"-"`
	City     string `ini:"city" arg:"-"`
	Region   string `ini:"region" arg:"-"`
	Postal   string `ini:"postal" arg:"-"`
	Country  string `ini:"country" arg:"-"`
	IpToken  string `ini:"token" arg:"-"`

	// other config info sent to page
	Realm    string `ini:"-" arg:"-"`
	Digests  string `ini:"-" arg:"-"`
	Admin    string `ini:"-" arg:"-"`
	Pass     string `ini:"-" arg:"-"`
	PublicIp string `ini:"-" arg:"-"`
}

type Weather struct {
	Timezone string `ini:"timezone"`
	Temps    string `ini:"temp"`
	Speeds   string `ini:"speed"`
	Depths   string `ini:"depth"`
}

var (
	// binding configs
	workingDir = "/var/lib/coventry"
	appDataDir = "/usr/share/apollo"
	etcPrefix  = "/etc"
	logPrefix  = "/var/log"
	version    = "unknown"
	publicIp   = "auto"

	// globals
	config  *Config  = nil
	weather *Weather = nil
	lock    sync.RWMutex
)

func (Config) Version() string {
	return "Version: " + version
}

func (Config) Description() string {
	return "apollo - web services for coventry phone system"
}

func init() {
	// parse arguments
	for pos, arg := range os.Args {
		switch arg {
		case "--":
			return
		case "-v":
			os.Args[pos] = "--verbose=1"
		case "-vv":
			os.Args[pos] = "--verbose=2"
		case "-vvv":
			os.Args[pos] = "--verbose=3"
		case "-vvvv":
			os.Args[pos] = "--verbose=4"
		case "-vvvvv":
			os.Args[pos] = "--verbose=5"
		}
	}
}

func load() {
	// default config
	new_config := Config{
		// service config
		Port:    8080,
		Keyfile: "./server.key",
		Crtfile: "./server.crt",

		// page config
		Home:     "https://www.tychosoft.com",
		Logo:     "/assets/logo.png",
		Views:    "en",
		Theme:    "dark",
		Location: "unspecified",
		PublicIp: publicIp,
		IpToken:  "*Your API Token*",
		Where:    "unknown",
		City:     "unknown",
		Region:   "unknown",
		Postal:   "unknown",
		Country:  "us",
	}

	new_weather := Weather{
		Timezone: "unknown",
		Temps:    "fahrenheit",
		Speeds:   "mph",
		Depths:   "inch",
	}

	configs, err := ini.LoadSources(ini.LoadOptions{Loose: true, Insensitive: true}, etcPrefix+"/apollo.conf", workingDir+"/custom.conf")
	if err == nil {
		configs.MapTo(&new_config)
		configs.Section("server").MapTo(&new_config)
		configs.Section("certs").MapTo(&new_config)
		configs.Section("weather").MapTo(&new_weather)
	} else {
		service.Error(err)
	}

	arg.MustParse(&new_config)
	if new_config.Host == "*" {
		new_config.Host = ""
	}

	err = apollo.Config(etcPrefix, workingDir)
	if err != nil {
		service.Error(err)
	}

	// set page values from full config...
	server := apollo.GetServer()
	forecast := apollo.GetWeather()
	new_config.Digests = apollo.Algorithm
	new_config.Realm = apollo.Realm
	new_config.Admin = apollo.GetConfig(server, "webadmin", "admin")
	new_config.Theme = apollo.GetConfig(server, "theme", "dark")
	new_config.IpToken = apollo.GetConfig(server, "token", "*Your API Token*")
	new_config.Location = apollo.GetConfig(server, "location", "unspecified")
	new_config.Where = apollo.GetConfig(server, "where", "none")
	new_config.City = apollo.GetConfig(server, "city", "unknown")
	new_config.Region = apollo.GetConfig(server, "region", "unknown")
	new_config.Postal = apollo.GetConfig(server, "postal", "unknown")

	new_weather.Timezone = apollo.GetConfig(forecast, "timezone", "unknown")

	common := apollo.GetCommon()
	new_config.Pass = apollo.GetConfig(common, "password", "")

	if dynCoventry == nil {
		dynInit(new_config.Port, new_config.Secure)
	}

	lock.Lock()
	defer lock.Unlock()
	config = &new_config
	weather = &new_weather
}

func main() {
	// config and setup service
	err := os.Chdir(workingDir)
	if err != nil {
		fmt.Println("Fatal: ", err)
		os.Exit(1)
	}

	load()
	service.Logger(config.Verbose, logPrefix+"/apollo.log")

	// setup app and routes
	address := fmt.Sprintf("%s:%v", config.Host, config.Port)
	aging := 600
	service.Debug(3, "prefix=", workingDir, ", bind=", address)
	service.Info("realm ", apollo.Realm, ", algo ", apollo.Algorithm)
	views := appDataDir + "/views_" + config.Views
	if flag, _ := apollo.IsDir(views); !flag {
		views = appDataDir + "/views"
	}
	if service.IsDebug() {
		aging = 10
	}

	engine := html.New(views, ".html")
	engine.Reload(service.IsDebug())
	app := fiber.New(fiber.Config{
		ServerHeader:          "Apollo",
		AppName:               "Apollo v" + version,
		DisableStartupMessage: true,
		Views:                 engine,
	})

	admin := basicauth.New(basicauth.Config{
		Authorizer: func(user, pass string) bool {
			if user != adminUser.Username {
				return false
			}

			digest := sha256.New()
			digest.Write([]byte(pass + ":" + user))
			return hex.EncodeToString(digest.Sum(nil)) == adminUser.Password
		},
		Realm: apollo.Realm,
	})

	user := func(ctx *fiber.Ctx) error {
		header := ctx.Get("Authorization")
		if len(header) < 8 || header[:7] != "Bearer " {
			return fiber.ErrUnauthorized
		}
		token := header[7:]
		id := apollo.VerifyToken(token)
		if id == 0 {
			return fiber.ErrUnauthorized
		}
		ctx.Locals("userID", id)
		return ctx.Next()
	}

	app.Static("/assets", appDataDir+"/assets", fiber.Static{MaxAge: aging})
	app.Post("/setup", postSetup)
	app.Post("/lines", admin, postNewLine)
	app.Post("/lines/:id", admin, postLine)
	app.Post("/lines/:id/delete", admin, deleteLine)
	app.Post("/lines/:id/passwd", admin, passwdLine)
	app.Post("/settings/theme", admin, themeSetup)
	app.Post("/settings/internet", admin, internetSetup)
	app.Post("/settings/location", admin, locationSetup)
	app.Delete("/lines/:id", admin, deleteLine)

	// client access api
	app.Get("/client/v0/ping", user, clientPing)
	app.Get("/client/v0/profile", user, clientProfile)
	app.Get("/client/v0/roster", user, clientRoster)
	app.Get("/client/v0/groups", user, clientGroups)

	// main views
	app.Get("/ping", admin, viewPing)
	app.Get("/main", admin, viewMain)
	app.Get("/lines", admin, viewLines)
	app.Get("/lines/:id", admin, editLine)
	app.Get("/groups", admin, viewGroups)
	app.Get("/contacts", admin, viewContacts)
	app.Get("/settings", admin, editSettings)
	app.Get("/setup", viewSetup)
	app.Get("/", func(ctx *fiber.Ctx) error {
		if setupFlag {
			return ctx.Redirect("/lines", fiber.StatusTemporaryRedirect)
		}
		return ctx.Redirect("/setup", fiber.StatusTemporaryRedirect)
	})

	// signal handler...
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		service.Live("start service on ", address)
		defer app.Shutdown()
		defer service.Stop("stop service")
		for {
			switch <-signals {
			case os.Interrupt: // sigint/ctrl-c
				fmt.Println()
				return
			case syscall.SIGTERM: // normal exit
				return
			case syscall.SIGHUP: // cleanup
				service.Reload("reload service")
				service.LoggerRestart()
				runtime.GC()
				load()
				service.Live()
			}
		}
	}()

	// start service(s)...
	if config.Secure {
		if err := app.ListenTLS(address, config.Crtfile, config.Keyfile); err != nil {
			service.Fail(99, err)
		}
	} else if err := app.Listen(address); err != nil {
		service.Fail(99, err)
	}
}
