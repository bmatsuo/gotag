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

	root := opt.Root
	force := opt.Force
	verbose := opt.Verbose

	gover, gorev, err := GetGoVersion()
	gotag := GoRepositoryTag(gover)
	if verbose {
		log.Printf("  Linker: %s", GoLinker)
		log.Printf(" Version: %s", gover)
		log.Printf("Revision: %d", gorev)
		log.Printf("     Tag: %s", gotag)
	}

	var project GoProject
	project, err = NewProject(root)
	Must(err)
	Must(BuildAndClean(project))

	var git Repository
	git, err = NewGitRepo(root)
	Must(err)

	if opt.Commit == "" { // It's OK to tag past commits if the HEAD is dirty.
		clean, err := git.IsClean()
		Must(err)
		if !clean {
			fmt.Fprint(os.Stderr, "The repository has uncommited changes.\n")
			fmt.Fprint(os.Stderr, "Commit the changes and run Gotag again.\n")
			os.Exit(1)
		}
	}

	if opt.Fetch {
		fmt.Print("Fetching remote tags\n")
		Must(git.TagsFetch())
	}

	var tags []string
	tags, err = git.Tags()
	Must(err)

	// Look for a tag named for the current version.
	hasCurrentTag := false
	for i := range tags {
		hasCurrentTag = gotag == tags[i]
		if hasCurrentTag {
			break
		}
	}

	// If found a try to delete it.
	if hasCurrentTag {
		fmt.Printf("Found tag %s\n", gotag)
		if force {
			Must(git.TagDelete(gotag))
		} else {
			fmt.Fprintf(os.Stderr, "use -f flag to update %s\n", gotag)
			os.Exit(1)
		}
	}

	// Create the new tag.
	annotation := fmt.Sprintf("Latest build for Go version %s %d", gover, gorev)
	if opt.Commit != "" {
		fmt.Fprintf(os.Stderr, "Creating tag %s %#v (%s)\n", gotag, annotation, opt.Commit)
		Must(git.TagNew(gotag, annotation, opt.Commit))
	} else {
		fmt.Fprintf(os.Stderr, "Creating tag %s %#v\n", gotag, annotation)
		Must(git.TagNew(gotag, annotation))
	}

	if opt.Push {
		fmt.Fprintf(os.Stderr, "Pushing tags to remote repository\n")
		Must(git.TagsPush())
	}
}
