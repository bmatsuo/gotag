// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    git.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Nov  6 01:17:31 PDT 2011
 *  Description: 
 */

import (
	"github.com/bmatsuo/go-script/script"
	"path/filepath"
	"errors"
	"strings"
	"fmt"
	"os"
)

type ErrorRepoConfig error

func NewErrorRepoConfig(err string) ErrorRepoConfig { return ErrorRepoConfig(errors.New(err)) }

type gitRepo struct {
	root  string
	shell script.Scriptor
}

func NewGitRepo(root string) (Repository, error) {
	repo := new(gitRepo)
	repo.root = root
	repo.shell = script.Bash
	return repo, repo.checkRoot()
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

func (repo *gitRepo) Root() string { return repo.root }
func (repo *gitRepo) Type() string { return "Git" }

func (repo *gitRepo) Name() (string, error) {
	// TODO: look at the contents of the .git/config file
	abs, err := filepath.Abs(repo.root)
	if err != nil {
		return "", err
	}
	return filepath.Base(abs), nil
}

func (repo *gitRepo) Tags() ([]string, error) {
	tagcmd := `git tag -l`
	if repo.root != "." {
		tagcmd = fmt.Sprintf("cd %s\n%s", script.ShellQuote(repo.root), tagcmd)
	}
	tagout, _, errexec := script.Output(repo.shell.NewScript(tagcmd))
	if errexec != nil {
		return nil, errexec
	}
	tags := strings.Fields(strings.Trim(string(tagout), "\n"))
	return tags, nil
}

func (repo *gitRepo) TagDelete(tag string) error {
	tagcmd := fmt.Sprintf(`git tag -d %s`, tag)
	if repo.root != "." {
		tagcmd = fmt.Sprintf("cd %s\n%s", script.ShellQuote(repo.root), tagcmd)
	}
	_, errexec := repo.shell.NewScript(tagcmd).Execute()
	return errexec
}

// If there is an extra value, it is used as a tag annotation.
// If there are two extra values, the second is used as a commit hash.
func (repo *gitRepo) Tag(name string, extra ...interface{}) error {
	tagcmd := fmt.Sprintf(`git tag`)
	if len(extra) > 0 {
		note := extra[0]
		switch note.(type) {
		case string:
			tagcmd = fmt.Sprintf(`%s -a -m %s`, tagcmd, script.ShellQuote(note.(string)))
		default:
			return errors.New("expected string annotation")
		}
	}
	tagcmd = fmt.Sprintf(`%s %s`, tagcmd, name)
	if len(extra) > 1 {
		tagcmd = fmt.Sprintf(`%s %s`, tagcmd, script.ShellQuote(extra[0].(string)))
	}
	if repo.root != "." {
		tagcmd = fmt.Sprintf("cd %s\n%s", script.ShellQuote(repo.root), tagcmd)
	}
	_, errexec := repo.shell.NewScript(tagcmd).Execute()
	return errexec
}
