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
	_ "embed"
	"fmt"

	"github.com/urfave/cli/v2"
)

var licenseCmd = &cli.Command{
	Hidden: true,
	Name:   "license",
	Usage:  "view the license header",
	Action: license,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "full",
			Aliases:  []string{"f"},
			Usage:    "view the full license instead of just thea header",
			Required: false,
		},
	},
}

func license(ctx *cli.Context) error {
	if ctx.Bool("full") {
		fmt.Println(licenseFullHeader + licenseFull)
	} else {
		fmt.Println(licenseHeader)
	}
	return nil
}

const licenseSmallHeader = `lwnfeed  Copyright (C) 2020-2022  Tulir Asokan

   This program comes with ABSOLUTELY NO WARRANTY.
   This is free software, and you are welcome to redistribute it
   under certain conditions; type ` + "`lwnfeed license`" + ` for details.`

const licenseHeader = `lwnfeed - A full-text RSS feed generator for LWN.net.
Copyright (C) 2020-2022 Tulir Asokan

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`

//go:embed LICENSE
var licenseFull string

const licenseFullHeader = `              lwnfeed - A full-text RSS feed generator for LWN.net.
                     Copyright (C) 2020-2022 Tulir Asokan

`
