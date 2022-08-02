// lwnfeed - A full-text RSS feed generator for LWN.net.
// Copyright (C) 2020-2022 Tulir Asokan
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
	"net/http"
	"strings"
	"time"

	log "maunium.net/go/maulogger/v2"
)

func writeError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(http.StatusText(status)))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}
	if time.Since(feed.Updated) > maxUpdateInterval {
		log.Debugln("Cached feed is over 30 minutes old, starting update")
		err := updateFeed()
		if err != nil {
			log.Errorln("Failed to update feed:", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
	}
	parts := strings.Split(r.URL.Path, "/")
	fileName := strings.ToLower(parts[len(parts)-1])
	var err error
	if fileName == "feed.rss" {
		err = feed.WriteRss(w)
	} else if fileName == "feed.atom" {
		err = feed.WriteAtom(w)
	} else if fileName == "feed.json" {
		err = feed.WriteJSON(w)
	} else {
		writeError(w, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Debugfln("Failed to write response to %s: %v", r.RemoteAddr, err)
	}
}

func serve(address string) error {
	log.Infoln("Listening on", address)
	return http.ListenAndServe(address, http.HandlerFunc(handleRequest))
}
