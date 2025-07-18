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
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func SHA256Secret(username, realm, password string) string {
	digest := sha256.New()
	digest.Write([]byte(username + ":" + realm + ":" + password))
	return hex.EncodeToString(digest.Sum(nil))
}

func MD5Secret(username, realm, password string) string {
	digest := md5.New()
	digest.Write([]byte(username + ":" + realm + ":" + password))
	return hex.EncodeToString(digest.Sum(nil))
}
