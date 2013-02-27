// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package model defines the basic
	data structures of the docs engine.
*/
package indexer

import (
	"andyk/docs/util"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Item struct {
	Path         string
	RenderedPath string
	Files        []File
	ChildItems   []Item
}

// Create a new repository item
func NewItem(path string, files []File, childItems []Item) Item {
	return Item{
		Path:         path,
		RenderedPath: getRenderedItemPath(path),
		Files:        files,
		ChildItems:   childItems,
	}
}

func (item Item) GetFilename() string {
	return filepath.Base(item.Path)
}

func (item *Item) GetHash() string {
	itemBytes, readFileErr := ioutil.ReadFile(item.Path)
	if readFileErr != nil {
		return ""
	}

	sha1 := sha1.New()
	sha1.Write(itemBytes)

	return fmt.Sprintf("%x", string(sha1.Sum(nil)[0:6]))
}

func (item Item) IsRendered() bool {
	return util.FileExists(item.RenderedPath)
}

func (item Item) GetRelativeItemPaths(basePath string) []string {

	paths := make([]string, 0, 0)

	for _, child := range item.ChildItems {
		paths = append(paths, child.GetRelativeItemPaths(basePath)...)
	}

	fullItemPath := item.RenderedPath
	relativePath := strings.Replace(fullItemPath, basePath, "", 1)
	paths = append(paths, relativePath)

	return paths
}

// Get the filepath of the rendered repository item
func getRenderedItemPath(itemPath string) string {
	itemDirectory := filepath.Dir(itemPath)
	itemName := strings.Replace(filepath.Base(itemPath), filepath.Ext(itemPath), "", 1)

	renderedFilePath := filepath.Join(itemDirectory, itemName+".html")
	return renderedFilePath
}
