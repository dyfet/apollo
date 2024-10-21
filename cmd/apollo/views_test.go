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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func TestViewPing(t *testing.T) {
	MockConfig()
	engine := html.New("../../web/views", ".html")
	engine.Reload(true)
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/ping", viewPing)

	req, _ := http.NewRequest("GET", "/ping", nil)
	resp, _ := app.Test(httptest.NewRequest(req.Method, req.URL.String(), nil))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := "<h1>Ping!</h1>"
	if !strings.Contains(string(body), expected) {
		t.Errorf("Expected to contain %q, but got %q", expected, string(body))
	}
}

func TestViewMain(t *testing.T) {
	MockConfig()
	engine := html.New("../../web/views", ".html")
	engine.Reload(true)
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/", viewMain)

	req, _ := http.NewRequest("GET", "/", nil)
	resp, _ := app.Test(httptest.NewRequest(req.Method, req.URL.String(), nil))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := "<h1>Home Page</h1>"
	if !strings.Contains(string(body), expected) {
		t.Errorf("Expected to contain %q, but got %q", expected, string(body))
	}
}
