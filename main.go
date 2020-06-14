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
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	Version   = "0.1.0"
	BuildTime string

	parsedBuildTime  time.Time
	workingDirectory string
)

var (
	rootURL  *url.URL
	loginURL *url.URL
	feedURL  *url.URL
)
var client http.Client

func main() {
	cli.VersionFlag.(*cli.BoolFlag).Aliases = []string{"V"}
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("lwnfeed v%s, built at %s\n", Version, parsedBuildTime.Format(time.RFC1123))
	}
	err := (&cli.App{
		Name:    "lwnfeed",
		Version: Version,
		Usage:   "A full-text RSS feed generator for LWN.net.",

		Commands: []*cli.Command{loginCmd, startCmd, licenseCmd},
		Flags: []cli.Flag{&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "the file to store the auth cookie in",
			Value:   filepath.Join(workingDirectory, "lwnfeed.cookie.gob"),
		}},

		Compiled: parsedBuildTime,
		Copyright: `lwnfeed  Copyright (C) 2020  Tulir Asokan

   This program comes with ABSOLUTELY NO WARRANTY.
   This is free software, and you are welcome to redistribute it
   under certain conditions; type ` + "`lwnfeed license`" + ` for details.`,

		HideHelpCommand:        true,
		UseShortOptionHandling: true,
	}).Run(os.Args)
	if err != nil {
		fail("", err)
	}
}
