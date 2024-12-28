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

package apollo

import (
	"testing"
)

func TestMD5Secret(t *testing.T) {
	username := "testuser"
	realm := "testrealm"
	password := "testpassword"
	expected := "3bfe9ca32df3fc2eb68848bb77c9d922"
	actual := MD5Secret(username, realm, password)
	if actual != expected {
		t.Errorf("Expected MD5Secret to return %q, but got %q", expected, actual)
	}
}

func TestSHA256Secret(t *testing.T) {
	username := "testuser"
	realm := "testrealm"
	password := "testpassword"
	expected := "2800485bba9476c3d949dbb360a70ec6dfe57a1e9075e4669193cfaaf2f08361"
	actual := SHA256Secret(username, realm, password)
	if actual != expected {
		t.Errorf("Expected SHA256Secret to return %q, but got %q", expected, actual)
	}
}
