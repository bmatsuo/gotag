// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    project.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Nov  6 01:06:52 PST 2011
 *  Description: 
 */

import (
	"github.com/bmatsuo/go-script/script"
	"path/filepath"
	"errors"
	"fmt"
	"os"
)

type GoProject interface {
	BuildType() string
	Root() string
	Name() (string, error)
	Build() error
	Clean() error
	Nuke() error
	Test() error
}

func BuildAndClean(gp GoProject) error {
	err := gp.Build()
	if err != nil {
		return err
	}
	return gp.Clean()
}

type Project struct {
	root string
	shell script.Scriptor
}

func NewProject(root string) (p *Project, err error) {
	p = new(Project)
	p.root = root
	p.shell = script.Bash
	return p, p.checkRoot()
}

func (p *Project) checkRoot() error {
	dir, err := os.Stat(p.root)
	if err != nil {
		return err
	}
	if !dir.IsDirectory() {
		return fmt.Errorf("Project root %s is not a directory.")
	}
	hasmakefile := false
	hasgofiles := false
	errwalk := filepath.Walk(p.root, func(path string, info *os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path != p.root && info.IsDirectory() {
			return filepath.SkipDir
		}
		if !hasmakefile && filepath.Base(path) == "Makefile" {
			hasmakefile = true
		} else if !hasgofiles && filepath.Ext(path) == ".go" {
			hasgofiles = true
		}
		return nil
	})
	if errwalk != nil {
		return errwalk
	}
	if !hasgofiles {
		return errors.New("no .go files")
	}
	if !hasmakefile {
		return errors.New("no Makefile")
	}
	return nil
}

func (p *Project) BuildType() string { return "gomake" }
func (p *Project) Root() string      { return p.root }
func (p *Project) Name() (string, error) {
	// TODO - Parse the Makefile for the project name.
	abs, err := filepath.Abs(p.root)
	if err != nil {
		return "", err
	}
	return filepath.Base(abs), nil
}

func (p *Project) Build() error {
	tagcmd := `gomake`
	if p.root != "." {
		tagcmd = fmt.Sprintf("cd %s\n%s", script.ShellQuote(p.root), tagcmd)
	}
	_, err := p.shell.NewScript(tagcmd).Execute()
	return err
}

func (p *Project) Clean() error {
	tagcmd := `gomake clean`
	if p.root != "." {
		tagcmd = fmt.Sprintf("cd %s\n%s", script.ShellQuote(p.root), tagcmd)
	}
	_, err := p.shell.NewScript(tagcmd).Execute()
	return err
}

func (p *Project) Nuke() error {
	tagcmd := `gomake nuke`
	if p.root != "." {
		tagcmd = fmt.Sprintf("cd %s\n%s", script.ShellQuote(p.root), tagcmd)
	}
	_, err := p.shell.NewScript(tagcmd).Execute()
	return err
}

func (p *Project) Test() error {
	tagcmd := `gotest`
	if p.root != "." {
		tagcmd = fmt.Sprintf("cd %s\n%s", script.ShellQuote(p.root), tagcmd)
	}
	_, err := p.shell.NewScript(tagcmd).Execute()
	return err
}
