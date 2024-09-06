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

	ipc "apollo/internal"
	"gitlab.com/tychosoft/service"
)

// Argument parser....
type Args struct {
	Config  string `arg:"--config" help:"server config file"`
	Host    string `arg:"env:APOLLO_HOST,--host" help:"server host address" default:""`
	Port    uint16 `arg:"env:APOLLO_PORT,--port" help:"server port" default:"8080"`
	Prefix  string `arg:"--prefix" help:"server path to coventry data"`
	Media   string `arg:"--media" help:"server path to bordeaux media"`
	Tls     bool   `arg:"--tls" help:"Server tls mode"`
	Verbose int    `arg:"-v" help:"debugging log level"`
}

// Server config
type Config struct {
	// apollo config data
	Prefix  string `ini:"prefix"`
	Media   string `ini:"media"`
	Host    string `ini:"host"`
	Port    uint16 `ini:"port"`
	Keyfile string `ini:"keyfile"`
	Crtfile string `ini:"crtfile"`
	Tls     bool   `ini:"tls"`

	// page specific and branding
	Logo     string `ini:"logo"`
	Home     string `ini:"home"`
	Views    string `ini:"views"`
	Theme    string `ini:"theme"`
	Weather  string `ini:"weather"`
	Location string `ini:"location"`
	Where    string `ini:"where"`
	City     string `ini:"city"`
	Region   string `ini:"region"`
	Postal   string `ini:"postal"`
	Country  string `ini:"country"`
	IpToken  string `ini:"token"`

	// other config info sent to page
	Realm    string `ini:"-"`
	Digests  string `ini:"-"`
	Admin    string `ini:"-"`
	Pass     string `ini:"-"`
	PublicIp string `ini:"-"`
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
	mediaData  = "/var/lib/bordeaux"
	etcPrefix  = "/etc"
	logPrefix  = "/var/log"
	version    = "unknown"
	publicIp   = "auto"

	// globals
	args    *Args    = &Args{Prefix: workingDir, Media: mediaData, Config: etcPrefix + "/apollo.conf"}
	config  *Config  = nil
	weather *Weather = nil
	lock    sync.RWMutex
)

func (Args) Version() string {
	return "Version: " + version
}

func (Args) Description() string {
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
	arg.MustParse(args)

	// config service
	logPath := logPrefix + "/apollo.log"
	service.Logger(args.Verbose, logPath)
	load()
	err := os.Chdir(config.Prefix)
	if err != nil {
		service.Fail(1, err)
	}
}

func load() {
	// default config
	new_config := Config{
		// service config
		Host:    args.Host,
		Port:    args.Port,
		Prefix:  args.Prefix,
		Tls:     args.Tls,
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

	configs, err := ini.LoadSources(ini.LoadOptions{Loose: true, Insensitive: true}, args.Config)
	if err == nil {
		// map and reset from args if not default
		configs.Section("server").MapTo(&new_config)
		configs.Section("weather").MapTo(&new_weather)
		if args.Port != 8080 {
			new_config.Port = args.Port
		}

		if len(args.Host) > 0 {
			new_config.Host = args.Host
		}

		if args.Prefix != workingDir {
			new_config.Prefix = args.Prefix
		}

		if args.Media != mediaData {
			new_config.Media = args.Media
		}

		if new_config.Host == "*" {
			new_config.Host = ""
		}
	} else {
		service.Error(err)
	}

	err = ipc.Config(etcPrefix, workingDir)
	if err != nil {
		service.Error(err)
	}

	// set page values from full config...
	server := ipc.GetServer()
	forecast := ipc.GetWeather()
	new_config.Digests = ipc.Algorithm
	new_config.Realm = ipc.Realm
	new_config.Admin = ipc.GetConfig(server, "webadmin", "admin")
	new_config.Theme = ipc.GetConfig(server, "theme", "dark")
	new_config.IpToken = ipc.GetConfig(server, "token", "*Your API Token*")
	new_config.Location = ipc.GetConfig(server, "location", "unspecified")
	new_config.Where = ipc.GetConfig(server, "where", "none")
	new_config.City = ipc.GetConfig(server, "city", "unknown")
	new_config.Region = ipc.GetConfig(server, "region", "unknown")
	new_config.Postal = ipc.GetConfig(server, "postal", "unknown")

	new_weather.Timezone = ipc.GetConfig(forecast, "timezone", "unknown")

	common := ipc.GetCommon()
	new_config.Pass = ipc.GetConfig(common, "password", "")

	if dynCoventry == nil {
		dynInit(new_config.Port, new_config.Tls)
	}

	lock.Lock()
	defer lock.Unlock()
	config = &new_config
	weather = &new_weather
}

func main() {
	// setup app and routes
	address := fmt.Sprintf("%s:%v", config.Host, config.Port)
	aging := 600
	service.Debug(3, "prefix=", config.Prefix, ", bind=", address)
	service.Info("realm ", ipc.Realm, ", algo ", ipc.Algorithm)
	views := appDataDir + "/views_" + config.Views
	if flag, _ := ipc.IsDir(views); !flag {
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
		Realm: ipc.Realm,
	})

	user := func(ctx *fiber.Ctx) error {
		header := ctx.Get("Authorization")
		if len(header) < 8 || header[:7] != "Bearer " {
			return fiber.ErrUnauthorized
		}
		token := header[7:]
		id := ipc.VerifyToken(token)
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
	if config.Tls {
		if err := app.ListenTLS(address, config.Crtfile, config.Keyfile); err != nil {
			service.Fail(99, err)
		}
	} else if err := app.Listen(address); err != nil {
		service.Fail(99, err)
	}
}
