// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    repository.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sat Nov  5 19:48:23 PDT 2011
 *  Description: 
 */

import (
	"fmt"
	"os"
)

type Repository interface {
	Root() string
	Tags() []string
}

type gitRepo struct {
	root  string
	shell Scriptor
}

func NewGitRepo(root string) (repo *gitRepo, err error) {
	repo = new(gitRepo)
	repo.root = root
	err = repo.checkRoot()
	return
}

func (repo *gitRepo) checkRoot() error {
	repo_dir := repo.root + "/.git"
	dir, staterr := os.Stat(repo_dir)
	if staterr != nil {
		return staterr
	}
	if !dir.IsDirectory() {
		return fmt.Errorf("Git file %s is not a directory.")
	}
	return nil
}
