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
	"strings"
    "log"
    "fmt"
    "os"
)

var opt options

func main() {
    opt = parseFlags()
	git, err := NewGitRepo(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	tags, err := git.Tags()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Join(tags, ":"))
}
