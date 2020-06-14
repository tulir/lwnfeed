// lwnfeed - A full-text RSS feed generator for LWN.net.
// Copyright (C) 2020 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/pkg/errors"
	log "maunium.net/go/maulogger/v2"
)

func saveCookies(cookie *http.Cookie, path string) error {
	fileHandle, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to open cookie file")
	}
	err = gob.NewEncoder(fileHandle).Encode(cookie)
	if err != nil {
		_ = fileHandle.Close()
		return errors.Wrap(err, "failed to write cookie to file")
	}
	err = fileHandle.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close cookie file")
	}
	return nil
}

func loadCookies(path string) error {
	log.Debugln("Reading cookies from", path)
	fileHandle, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "failed to open cookie file")
	}
	var cookie http.Cookie
	err = gob.NewDecoder(fileHandle).Decode(&cookie)
	if err != nil {
		return errors.Wrap(err, "failed to read cookie from file")
	}
	client.Jar.SetCookies(rootURL, []*http.Cookie{&cookie})
	_ = fileHandle.Close()
	return nil
}
