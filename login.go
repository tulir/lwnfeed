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

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var loginCmd = &cli.Command{
	Name:   "login",
	Usage:  "log into LWN and store the auth cookie.",
	Action: login,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Usage:   "username to log in as (will prompt if omitted)",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Usage:   "password to log in with (will prompt if omitted)",
		},
	},
}

var empty = errors.New("value is empty")

func notEmpty(val string) error {
	if len(val) == 0 {
		return empty
	}
	return nil
}

func noRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

func login(ctx *cli.Context) error {
	username := ctx.String("username")
	password := ctx.String("password")
	var err error
	if len(username) == 0 {
		username, err = (&promptui.Prompt{Label: "Username", Validate: notEmpty}).Run()
		if err != nil {
			return err
		}
	}
	if len(password) == 0 {
		password, err = (&promptui.Prompt{Label: "Password", Mask: '*', Validate: notEmpty}).Run()
		if err != nil {
			return err
		}
	}
	client.CheckRedirect = noRedirect
	resp, err := client.PostForm(loginURL.String(), url.Values{
		"Username": []string{username},
		"Password": []string{password},
		"target":   []string{""},
		"submit":   []string{"Log+in"},
	})
	client.CheckRedirect = nil
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to log in: login responded %s", resp.Status)
	}

	var authCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "LWNSession1" {
			authCookie = cookie
		}
	}
	if authCookie == nil {
		return fmt.Errorf("failed to log in: login response did not contain LWNSession1 cookie")
	}

	client.Jar.SetCookies(rootURL, []*http.Cookie{authCookie})

	if resp.StatusCode >= 300 {
		resp, err = client.Get(resp.Header.Get("Location"))
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("failed to log in: login redirect responded %s", resp.Status)
		}
	} else {
		fmt.Println(resp.Status)
		fmt.Println("Warning: login response didn't respond with redirect")
	}

	err = saveCookies(authCookie, ctx.Path("file"))
	if err != nil {
		return err
	}

	fmt.Println("Successfully wrote auth cookies to file")
	return nil
}
