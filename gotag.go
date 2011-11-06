// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    gotag.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sat Nov  5 19:46:28 PDT 2011
 *  Description: Main source file in gotag
 */

import (
	"github.com/bmatsuo/go-script/script"
	"runtime"
	"strconv"
	"strings"
	"log"
	"fmt"
	"os"
)

var archlinker = map[string]string{
	"amd64": "6l",
	"386":   "8l",
	"arm":   "5l",
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var GoLinker = archlinker[runtime.GOARCH]

func GetGoVersion() (version string, revision int, err error) {
	if GoLinker == "" {
		panic(fmt.Errorf("unknown architechture %s", runtime.GOARCH))
	}
	var p []byte
	p, _, err = script.Output(script.Bash.NewScript(fmt.Sprintf("%s -V", GoLinker)))
	if err != nil {
		return
	}
	pieces := strings.Fields(string(p))
	if len(pieces) < 2 {
		err = fmt.Errorf("Didn't understand Go version %s", string(p))
	}
	version = pieces[len(pieces)-2]
	revision, err = strconv.Atoi(pieces[len(pieces)-1])
	return
}

func GoRepositoryTag(version string) string { return "go." + version }

var opt options

func main() {
	opt = parseFlags()

	gover, gorev, err := GetGoVersion()
	log.Printf("  Linker: %s", GoLinker)
	log.Printf(" Version: %s", gover)
	log.Printf("Revision: %d", gorev)
	gotag := GoRepositoryTag(gover)
	log.Printf("     Tag: %s", gotag)

	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	var project GoProject
	project, err = NewProject(root)
	Must(err)
	Must(BuildAndClean(project))

	var git Repository
	git, err = NewGitRepo(root)
	Must(err)

	var tags []string
	tags, err = git.Tags()
	Must(err)

	hasCurrentTag := false
	for i := range tags {
		if gotag == tags[i] {
			hasCurrentTag = true
		}
	}
	dodelete := false
	if hasCurrentTag {
		fmt.Printf("Tag %s found. It must be deleted.\n", gotag)
		if dodelete {
			Must(git.TagDelete(gotag))
		}
	} else {
		log.Printf("No tag %s", gotag)
	}

	Must(git.Tag(gotag, fmt.Sprintf("Latest build for Go version %s %d", gover, gorev)))
	log.Printf("Tagged")
	tags, err = git.Tags()
	Must(err)
	fmt.Println(strings.Join(tags, ":"))
	Must(git.TagDelete(gotag))
	tags, err = git.Tags()
	Must(err)
	fmt.Println(strings.Join(tags, ":"))
}
