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
	"sync"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	log "maunium.net/go/maulogger/v2"
)

var feed = &feeds.Feed{
	Title:       "LWN.net",
	Link:        &feeds.Link{Href: "https://lwn.net"},
	Description: "LWN.net is a comprehensive source of news and opinions from and about the Linux community. This is the main LWN.net feed, listing all articles which are posted to the site front page.",
	Created:     time.Now(),
}
var feedLoadLock sync.Mutex
var maxUpdateInterval = 30 * time.Minute

func initFeed() {
	log.Debugln("Starting initial feed load")
	err := updateFeed()
	if err != nil {
		log.Errorln("Failed to load initial feed:", err)
	}
}

func updateFeed() error {
	feedLoadLock.Lock()
	defer feedLoadLock.Unlock()
	if time.Now().Sub(feed.Updated) < maxUpdateInterval {
		return nil
	}

	log.Infoln("Updating feed from LWN.net...")
	resp, err := client.Get(feedURL.String())
	if err != nil {
		return errors.Wrap(err, "failed to fetch feed")
	}
	defer resp.Body.Close()
	inputFeed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to parse feed")
	}
	log.Debugln("Received feed with", len(inputFeed.Items), "items, fetching items...")

	outputItems := make([]*feeds.Item, len(inputFeed.Items))
	var ptr int
	for _, inputItem := range inputFeed.Items {
		outputItem, err := handleInputFeedItem(inputItem)
		if err != nil {
			log.Warnfln("Failed to handle item %s: %v", inputItem.Link, err)
			continue
		}
		outputItems[ptr] = outputItem
		ptr++
	}
	feed.Items = outputItems[:ptr]
	feed.Updated = time.Now()
	log.Infoln("Feed updated with", len(feed.Items), "items")
	return nil
}
