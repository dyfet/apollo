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
	"testing"
)

func TestHasMD5(t *testing.T) {
	Algorithm = "MD5"
	if !HasMD5() {
		t.Errorf("Expected HasMD5 to return true for Algorithm = MD5")
	}
	Algorithm = "SHA-256"
	if HasMD5() {
		t.Errorf("Expected HasMD5 to return false for Algorithm = SHA-256")
	}
}

func TestHasSHA256(t *testing.T) {
	Algorithm = "SHA-256"
	if !HasSHA256() {
		t.Errorf("Expected HasSHA256 to return true for Algorithm = SHA-256")
	}
	Algorithm = "MD5"
	if HasSHA256() {
		t.Errorf("Expected HasSHA256 to return false for Algorithm = MD5")
	}
}

func TestComputeMD5(t *testing.T) {
	id := "test-id"
	secret := "test-secret"
	Realm = "test-realm"
	expected := "9c5a7b075f77a2733b06551d835aa469"
	actual := ComputeMD5(id, secret)
	if actual != expected {
		t.Errorf("Expected ComputeMD5 to return %q, but got %q", expected, actual)
	}
}

func TestComputeSHA256(t *testing.T) {
	id := "test-id"
	secret := "test-secret"
	Realm = "test-realm"
	expected := "bdf1798febaf6c8b1c9bea93a98361e75a3833974060f311f936ba56ce8fd961"
	actual := ComputeSHA256(id, secret)
	if actual != expected {
		t.Errorf("Expected ComputeSHA256 to return %q, but got %q", expected, actual)
	}
}
