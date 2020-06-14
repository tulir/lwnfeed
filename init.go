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
	"fmt"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/publicsuffix"
)

func fail(message string, err error) {
	if len(message) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v", message, err)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
	}
	os.Exit(2)
}

func init() {
	var err error

	parsedBuildTime, err = time.Parse(time.RFC3339, BuildTime)
	if err != nil {
		fail("Failed to parse build time", err)
	}

	workingDirectory, err = os.Getwd()
	if err != nil {
		fail("Failed to get working directory", err)
	}

	rootURL, err = url.Parse("https://lwn.net")
	if err != nil {
		fail("Failed to parse URL", err)
	}
	loginURL, err = url.Parse("https://lwn.net/Login/")
	if err != nil {
		fail("Failed to parse URL", err)
	}
	feedURL, err = url.Parse("https://lwn.net/headlines/rss")
	if err != nil {
		fail("Failed to parse URL", err)
	}

	client.Jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		fail("Failed to create cookie jar", err)
	}
}
