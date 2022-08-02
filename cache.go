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
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/gorilla/feeds"
	log "maunium.net/go/maulogger/v2"
)

var cachedArticles = make(map[int]*feeds.Item)
var cacheFile *os.File
var cacheWriter *gob.Encoder

type CacheItem struct {
	ID   int
	Item *feeds.Item
}

func readCache(path string) error {
	log.Infoln("Loading cached articles from", path)
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	dec := gob.NewDecoder(file)
	for {
		var cacheItem CacheItem
		err = dec.Decode(&cacheItem)
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error decoding item: %w", err)
		}
		cachedArticles[cacheItem.ID] = cacheItem.Item
	}
	log.Infoln("Loaded", len(cachedArticles), "cached articles from disk")
	return nil
}

func openCacheWriter(path string) (err error) {
	log.Debugln("Opening cache file for writing")
	cacheFile, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	cacheWriter = gob.NewEncoder(cacheFile)
	log.Debugln("Rewriting cache file on disk")
	for id, item := range cachedArticles {
		err = cacheWriter.Encode(CacheItem{id, item})
		if err != nil {
			log.Warnfln("Failed to add article %d to disk cache: %v", id, err)
		}
	}
	log.Debugln("Successfully opened cache file")
	return
}

func addToCache(id int, item *feeds.Item) {
	cachedArticles[id] = item
	if cacheWriter != nil {
		err := cacheWriter.Encode(CacheItem{id, item})
		if err != nil {
			log.Warnfln("Failed to add article %d to disk cache: %v", id, err)
		}
	}
	log.Debugln("Added article", id, "to cache")
}
