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
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	log "maunium.net/go/maulogger/v2"
)

var (
	unknownLink   = errors.New("unknown feed item link format")
	invalidItemID = errors.New("feed item link has invalid format")
)

func loadArticleContent(id int) (*feeds.Item, error) {
	log.Debugfln("Loading content of article %d from LWN.net", id)
	link, err := url.Parse(fmt.Sprintf("https://lwn.net/Articles/%d/", id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse URL")
	}

	resp, err := client.Get(link.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch article")
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected response status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse article HTML")
	}

	title := doc.Find(".PageHeadline > h1").Text()
	body, err := doc.Find(".ArticleText").Html()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get article body")
	}

	log.Infofln("Successfully loaded content of article %d", id)

	return &feeds.Item{
		Link:    &feeds.Link{Href: link.String()},
		Title:   title,
		Content: body,
		Id:      strconv.Itoa(id),
	}, nil
}

func handleInputFeedItem(input *gofeed.Item) (*feeds.Item, error) {
	if !strings.HasPrefix(input.Link, "https://lwn.net/Articles/") {
		return nil, unknownLink
	}
	id, err := strconv.Atoi(strings.Split(input.Link, "/")[4])
	if err != nil {
		return nil, invalidItemID
	}
	output, ok := cachedArticles[id]
	if ok {
		log.Debugln("Found article", id, "in cache")
		return output, nil
	}
	output, err = loadArticleContent(id)
	if err != nil {
		return nil, err
	}
	output.Description = input.Description
	output.Created = *input.PublishedParsed
	output.Author = &feeds.Author{
		Name:  input.Author.Name,
		Email: input.Author.Email,
	}
	addToCache(id, output)
	return output, nil
}
