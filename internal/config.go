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

package internal

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/ini.v1"

	"gitlab.com/tychosoft/service"
)

type Line struct {
	Caller   string `ini:"caller" json:"caller"`
	Display  string `ini:"display" json:"display"`
	Lines    uint16 `ini:"lines" json:"lines"`
	Type     string `ini:"type" json:"type"`
	Location string `ini:"location" json:"location"`
	Cabling  string `ini:"cabling" json:"cabling"`
	EMail    string `ini:"email" json:"email"`
	Secret   string `ini:"secret" json:"-"`
	MD5      string `ini:"md5" json:"-"`
	SHA256   string `ini:"sha256" json:"-"`
	ACL      string `ini:"acl" json:"-"`
	COVERAGE string `ini:"coverage" json:"-"`
	DELAYED  string `ini:"delayed" json:"-"`
	Count    uint16 `ini:"-" json:"count"`
	Presence string `ini:"-" json:"status"`
	Agent    string `ini:"-" json:"agent"`
	Editable bool   `ini:"-" json:"-"`
	Host     string `ini:"-" json:"-"`
	URL      string `ini:"-" json:"-"`
}

type Group struct {
	Display string `json:"display"`
	Members []int  `json:"members"`
	Type    string `json:"-"`
}

var (
	// internal globals
	Realm     = ""
	Algorithm = "SHA-256"
	Password  = ""

	// local vars
	coventryConfig *ini.File = nil
	coventryUpdate *ini.File = nil
	coventryCustom *ini.File = nil
	coventrySaveTo string
	lock           sync.RWMutex
)

func UpdateCoventry(group, key, value string) error {
	var err error
	section := coventryUpdate.Section(group)
	if section == nil {
		section, err = coventryUpdate.NewSection(group)
		if err != nil {
			return err
		}
	}
	SetConfig(section, key, value)
	return nil
}

func SaveCoventry() error {
	err := coventryUpdate.SaveTo(coventrySaveTo)
	if err != nil {
		return err
	}

	return ReloadCoventry()
}

func defaultConfig() error {
	var err error = nil
	var section *ini.Section = nil

	section = coventryConfig.Section("server")
	if !section.HasKey("sitename") {
		_, err = section.NewKey("sitename", "Coventry Server")
	}

	if err == nil && !section.HasKey("location") {
		_, err = section.NewKey("location", "unspecified")
	}

	if err == nil && len(Realm) < 1 {
		Realm, err = os.Hostname()
		Realm = GetConfig(section, "hostname", Realm)
		Realm = GetConfig(section, "realm", Realm)
	}

	if err == nil {
		algo := strings.ToUpper(GetConfig(section, "algorithm", Algorithm))
		if strings.Contains(algo, "MD5") && strings.Contains(algo, "SHA") {
			Algorithm = "SHA-256, MD5"
		} else if strings.Contains(algo, "MD5") {
			Algorithm = "MD5"
		} else if strings.Contains(algo, "SHA") {
			Algorithm = "SHA-256"
		}
	}

	section = coventryConfig.Section("messages")
	if err == nil && !section.HasKey("welcome") {
		_, err = section.NewKey("welcome", "Welcome to Coventry")
	}
	if err == nil && !section.HasKey("shutdown") {
		_, err = section.NewKey("shutdown", "Shutting down...")
	}

	section = coventryConfig.Section("common")
	if err == nil && section.HasKey("password") {
		Password = GetConfig(section, "password", "")
	}
	if err == nil && !section.HasKey("lines") {
		_, err = section.NewKey("lines", "1")
	}
	if err == nil && !section.HasKey("presence") {
		_, err = section.NewKey("presence", "here")
	}
	if err == nil && !section.HasKey("type") {
		_, err = section.NewKey("type", "generic")
	}
	if err == nil && !section.HasKey("room") {
		_, err = section.NewKey("room", "any")
	}
	if err == nil && !section.HasKey("location") {
		_, err = section.NewKey("location", "unspecified")
	}

	section = coventryConfig.Section("calls")
	if err == nil && !section.HasKey("mode") {
		_, err = section.NewKey("mode", "proxy")
	}
	if err == nil && !section.HasKey("ring") {
		_, err = section.NewKey("ring", "4")
	}
	if err == nil && !section.HasKey("delayed") {
		_, err = section.NewKey("delayed", "12")
	}

	if err == nil && !coventryConfig.HasSection("features") {
		section = coventryConfig.Section("features")
		if !section.HasKey("*99") {
			_, err = section.NewKey("*99", "echo")
		}
		if err == nil && !section.HasKey("*98") {
			_, err = section.NewKey("*98", "reload")
		}
		if err == nil && !section.HasKey("*97") {
			_, err = section.NewKey("*97", "@weather")
		}
	}
	return err
}

func SetConfig(section *ini.Section, id string, value string) {
	key, err := section.GetKey(id)
	if err == nil {
		key.SetValue(value)
	} else {
		section.NewKey(id, value)
	}
}

func GetConfig(section *ini.Section, id string, def string) string {
	key, err := section.GetKey(id)
	if err == nil {
		return key.Value()
	}
	return def
}

func Config(etcPrefix string, covPrefix string) error {
	var err error
	var opt = ini.LoadOptions{Loose: true, Insensitive: true}

	lock.Lock()
	defer lock.Unlock()
	ipcInit(covPrefix + "/ipc.json")

	coventrySaveTo = covPrefix + "/dynamic.conf"
	coventryUpdate, _ = ini.Load(covPrefix + "/dynamic.conf")
	coventryCustom, _ = ini.Load(covPrefix + "/custom.conf")
	coventryConfig, err = ini.LoadSources(opt,
		etcPrefix+"/coventry.conf",
		covPrefix+"/dynamic.conf",
		covPrefix+"/custom.conf",
		// covPrefix+"/state.conf", - not computed for base config
	)

	if err != nil {
		service.Fail(99, err)
	}

	if coventryConfig == nil {
		coventryConfig = &ini.File{}
	}

	if coventryUpdate == nil {
		coventryUpdate = &ini.File{}
	}

	if err != nil {
		defaultConfig()
		return err
	}

	return defaultConfig()
}

func GetCommon() *ini.Section {
	lock.RLock()
	defer lock.RUnlock()
	return coventryConfig.Section("common")
}

func GetServer() *ini.Section {
	lock.RLock()
	defer lock.RUnlock()
	return coventryConfig.Section("server")
}

func GetWeather() *ini.Section {
	lock.RLock()
	defer lock.RUnlock()
	return coventryConfig.Section("weather")
}

func GetRooms() *ini.Section {
	lock.RLock()
	defer lock.RUnlock()
	return coventryConfig.Section("rooms")
}

func GetZones() *ini.Section {
	lock.RLock()
	defer lock.RUnlock()
	return coventryConfig.Section("zones")
}

func GetFeatures() *ini.Section {
	lock.RLock()
	defer lock.RUnlock()
	return coventryConfig.Section("features")
}

func GetGroups() map[string]*Group {
	groups := make(map[string]*Group)
	lock.RLock()
	defer lock.RUnlock()
	section := coventryConfig.Section("groups")
	for _, key := range section.Keys() {
		id := key.Name()
		group := fetchGroup(id)
		if group != nil {
			groups[id] = group
		}
	}
	return groups
}

func GetPolicies() map[string]*Group {
	groups := make(map[string]*Group)
	lock.RLock()
	defer lock.RUnlock()
	section := coventryConfig.Section("access")
	for _, key := range section.Keys() {
		id := key.Name()
		group := fetchPolicy(id)
		if group != nil {
			groups[id] = group
		}
	}

	section = coventryConfig.Section("groups")
	for _, key := range section.Keys() {
		id := key.Name()
		group := fetchGroup(id)
		if group != nil {
			groups[id] = group
		}
	}

	return groups
}

func GetLines() map[int]*Line {
	lines := make(map[int]*Line)
	lock.RLock()
	defer lock.RUnlock()

	for _, section := range coventryConfig.Sections() {
		key := section.Name()
		if key < "10" || key > "89" {
			continue
		}
		id, _ := strconv.Atoi(key)
		line := &Line{Agent: "offline", URL: "none", Presence: "down", Count: 0, Editable: true}
		coventryConfig.Section("common").MapTo(line)
		coventryConfig.Section(key).MapTo(line)
		getRegistry(id, line)
		sec := coventryCustom.Section(key)
		if len(sec.Keys()) > 0 {
			line.Editable = false
		}
		lines[id] = line
	}
	return lines
}

func GetGroup(id string) *Group {
	if id < "100" {
		return nil
	}

	lock.RLock()
	defer lock.RUnlock()
	return fetchGroup(id)
}

func GetPolicy(id string) *Group {
	if id < "100" {
		return nil
	}

	lock.RLock()
	defer lock.RUnlock()
	group := fetchGroup(id)
	if group != nil {
		return group
	}
	return fetchPolicy(id)
}

func fetchGroup(id string) *Group {
	group := &Group{}
	groups := coventryConfig.Section("groups")
	display := coventryConfig.Section("display")
	if !groups.HasKey(id) {
		if id == "system" {
			return group
		}
		return nil
	}

	group.Type = "group"
	if display.HasKey(id) {
		key, err := display.GetKey(id)
		if err == nil {
			group.Display = key.Value()
		}
	}

	key, err := groups.GetKey(id)
	if err != nil {
		return group
	}
	group.Members = getMembers(key.Value())
	return group
}

func fetchPolicy(id string) *Group {
	group := &Group{}
	access := coventryConfig.Section("access")
	if !access.HasKey(id) {
		return nil
	}

	group.Type = "access"
	key, err := access.GetKey(id)
	if err != nil {
		return group
	}
	group.Members = getMembers(key.Value())
	return group
}

func getMembers(members string) []int {
	var out []int
	list := strings.FieldsFunc(members, func(r rune) bool {
		return r == ',' || r == ';' || r == ':' || r == ' ' || r == '\t'
	})
	for _, str := range list {
		if len(str) < 1 {
			continue
		}
		member, err := strconv.Atoi(str)
		if err != nil {
			continue
		}
		if member < 10 || member > 89 {
			continue
		}
		mid := strconv.Itoa(member)
		if !coventryConfig.HasSection(mid) {
			continue
		}
		out = append(out, member)
	}
	return out
}

func updateKeys(cfg *ini.File, id string, data interface{}) error {
	section, err := cfg.NewSection(id)
	if err != nil {
		return err
	}

	v := reflect.Indirect(reflect.ValueOf(data))
	t := v.Type()

	for pos := 0; pos < v.NumField(); pos++ {
		field := t.Field(pos)
		tag := field.Tag.Get("ini")
		if tag != "" && tag != "-" {
			value := v.Field(pos).Interface()
			change := fmt.Sprintf("%v", value)
			if len(change) > 0 {
				section.Key(tag).SetValue(change)
			} else {
				section.DeleteKey(tag)
			}
		}
	}
	return nil
}

func UpdateLine(extension int, line *Line) error {
	if extension < 10 || extension > 89 {
		return fmt.Errorf("invalid extension number")
	}

	id := strconv.Itoa(extension)
	coventryUpdate.DeleteSection(id)
	err := updateKeys(coventryUpdate, id, line)
	if err == nil {
		err = coventryUpdate.SaveTo(coventrySaveTo)
	}

	if err != nil {
		return err
	}

	return ReloadCoventry()
}

func RemoveLine(extension int) error {
	if extension < 10 || extension > 89 {
		return fmt.Errorf("invalid extension number")
	}

	coventryUpdate.DeleteSection(strconv.Itoa(extension))

	err := coventryUpdate.SaveTo(coventrySaveTo)
	if err != nil {
		return err
	}

	return ReloadCoventry()
}

func ExistsLine(extension int) bool {
	id := strconv.Itoa(extension)
	for _, sec := range coventryConfig.Sections() {
		if sec.Name() == id {
			return true
		}
	}
	return false
}

func NewLine() (int, *Line) {
	line := &Line{Lines: 1, Type: "generic", Location: "unspecified", Editable: true}
	lock.RLock()
	defer lock.RUnlock()
	coventryConfig.Section("common").MapTo(line)

	for ext := 10; ext <= 89; ext++ {
		id := strconv.Itoa(ext)
		sec := coventryConfig.Section(id)
		if len(sec.Keys()) == 0 {
			coventryConfig.DeleteSection(id)
			return ext, line
		}
	}
	return 0, nil
}

func SavedLine(extension int) *Line {
	if extension < 10 || extension > 89 {
		return nil
	}

	line := &Line{}
	key := strconv.Itoa(extension)
	lock.RLock()
	defer lock.RUnlock()
	coventryUpdate.Section(key).MapTo(line)
	return line
}

func GetLine(extension int) *Line {
	if extension < 10 || extension > 89 {
		return nil
	}

	line := &Line{Agent: "offline", URL: "none", Presence: "down", Count: 0, Editable: true}
	key := strconv.Itoa(extension)
	lock.RLock()
	defer lock.RUnlock()
	coventryConfig.Section("common").MapTo(line)
	coventryConfig.Section(key).MapTo(line)
	getRegistry(extension, line)
	sec := coventryCustom.Section(key)
	if len(sec.Keys()) > 0 {
		line.Editable = false
	}

	return line
}

func CountLines() int {
	lines := 0
	lock.RLock()
	defer lock.RUnlock()

	for _, section := range coventryConfig.Sections() {
		key := section.Name()
		if key < "10" || key > "89" {
			continue
		}
		lines++
	}
	return lines
}
