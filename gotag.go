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
	"template"
	"runtime"
	"strconv"
	"strings"
	"bytes"
	"log"
	"fmt"
	"os"
)

var tfuncs = template.FuncMap{
	"quote": func(x interface{}) (string, error) {
		switch x.(type) {
		case string:
			return script.ShellQuote(x.(string)), nil
		}
		return "", fmt.Errorf("argument %#v is not a string", x)
	},
}

var cmdtemplates = `
{{/* Outputs a shell command given a list of strings (executable + args) */}}
	{{define "cmd"}}{{if ""}}
		{{end}}{{with $cmd := .}}{{range $i, $arg := $cmd}}{{if ""}}
			{{end}}{{if $i}} {{end}}{{quote $arg}}{{end}}{{end}}{{end}}

{{/* Outputs a list of comands .cmds. If .dir is set, the working directory is set with cd*/}}
	{{define "script"}}{{if ""}}
				{{end}}{{if .dir}}cd {{quote .dir}}
{{end}}{{if ""}}
				{{end}}{{range $i, $cmd := .cmds}}{{if $i}}
{{end}}{{if ""}}
				{{end}}{{template "cmd" $cmd}}{{end}}{{end}}
`

var templates = template.SetMust(new(template.Set).Funcs(tfuncs).Parse(cmdtemplates))

type ShellCmd []string

func CmdTemplateScript(sh script.Scriptor, dir string, cmds ...ShellCmd) script.Script {
	if sh == nil {
		panic("nil scriptor")
	}
	var d string
	if dir != "." {
		d = dir
	}
	buff := new(bytes.Buffer)
	err := templates.Template("script").Execute(buff, map[string]interface{}{"dir": d, "cmds": cmds})
	if err != nil {
		log.Println(err)
	}
	return sh.NewScript(string(buff.Bytes()))
}

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
	force := false
	if hasCurrentTag {
		fmt.Printf("found tag %s\n", gotag)
		if force {
			fmt.Printf("deleting %s\n", gotag)
			Must(git.TagDelete(gotag))
		} else {
			fmt.Fprintf(os.Stderr, "use -f flag to update %s\n", gotag)
			os.Exit(1)
		}
	}
	CmdTemplateScript(script.Bash, "some/path",
		ShellCmd{"echo", "hello, worlrd"},
		ShellCmd{"cd", "google"},
		ShellCmd{"goma'ke", "nuke"})

	annotation := fmt.Sprintf("Latest build for Go version %s %d", gover, gorev)
	fmt.Fprintf(os.Stderr, "creating tag %s %#v\n", gotag, annotation)
	Must(git.Tag(gotag, annotation))
	log.Printf("Tagged")
	tags, err = git.Tags()
	Must(err)
	fmt.Println(strings.Join(tags, ":"))
	Must(git.TagDelete(gotag))
	log.Printf("Deleted")
	tags, err = git.Tags()
	tags, err = git.Tags()
	Must(err)
	fmt.Println(strings.Join(tags, ":"))
}
