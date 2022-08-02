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
	"github.com/urfave/cli/v2"
	log "maunium.net/go/maulogger/v2"
)

var startCmd = &cli.Command{
	Name:   "start",
	Usage:  "start serving the feed.",
	Action: start,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "listen",
			Aliases: []string{"l"},
			Value:   "localhost:8080",
			Usage:   "the address to listen on",
		},
		&cli.PathFlag{
			Name:    "cache",
			Aliases: []string{"c"},
			Usage:   "the file to store cached article content in (optional)",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "log more debug information",
		},
	},
}

func start(ctx *cli.Context) error {
	if ctx.Bool("verbose") {
		log.DefaultLogger.PrintLevel = log.LevelDebug.Severity
	}
	err := loadCookies(ctx.Path("file"))
	if err != nil {
		return err
	}
	err = loadCache(ctx)
	if err != nil {
		return err
	}
	go initFeed()
	return serve(ctx.String("listen"))
}

func loadCache(ctx *cli.Context) error {
	if ctx.IsSet("cache") {
		path := ctx.Path("cache")
		err := readCache(path)
		if err != nil {
			return err
		}
		err = openCacheWriter(path)
		if err != nil {
			return err
		}
	}
	return nil
}
