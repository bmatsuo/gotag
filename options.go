// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    options.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sat Nov  5 19:46:28 PDT 2011
 *  Description: Option parsing for gotag
 */

import (
	"flag"
	"fmt"
	"os"
)
/*
 *  Constants, variables, and functions that users may actually want to call
 *  are capitalized.
 */

var (
	// Set this variable to customize the help message header.
	// For example, `gotag [options] action [arg2 ...]`.
	CommandLineHelpUsage string
	// Set this variable to print a message after the option specifications.
	// For example, "For more help:\n\tgotag help [action]"
	CommandLineHelpFooter string
)

//  A struct that holds parsed option values.
//  TODO: Customize this struct with options for gotag
type options struct {
	Root    string
	Fetch   bool
	Test    bool
	Push    bool
	Force   bool
	Commit  string
	Verbose bool
}

//  Create a flag.FlagSet to parse the command line options/arguments.
func setupFlags(opt *options) *flag.FlagSet {
	fs := flag.NewFlagSet("gotag", flag.ExitOnError)
	fs.StringVar(&opt.Commit, "commit", "", "Specify commit to tag.")
	fs.BoolVar(&opt.Fetch, "fetch", true, "Fetch remote tags before creating new tags.")
	fs.BoolVar(&opt.Test, "test", true, "TEMPORARY TEST FLAG.")
	fs.BoolVar(&opt.Push, "push", true, "Push newly created tags when finished.")
	fs.BoolVar(&opt.Force, "f", false, "Delete existing tag if conflict exists.")
	fs.BoolVar(&opt.Verbose, "v", false, "Verbose program output.")

	setupUsage(fs)
	return fs
}

//  Check the options for acceptable values. Panics or otherwise exits
//  with a non-zero exitcode when errors are encountered.
//  TODO: Make sure the gotag's flags are valid.
func verifyFlags(opt *options, fs *flag.FlagSet) {
	args := fs.Args()
	if len(args) < 1 {
		opt.Root = "."
	} else {
		opt.Root = args[0]
	}
	if info, err := os.Stat(opt.Root); err != nil {
		fs.Usage()
		fmt.Fprintf(os.Stderr, "stat error: %s\n", err.Error())
	} else if !info.IsDirectory() {
		fs.Usage()
		fmt.Fprintf(os.Stderr, "ROOT %s is not a directory\n", opt.Root)
	}
}

//  Print a help message to standard error. See constants CommandLineHelpUsage
//  and CommandLineHelpFooter.
func PrintHelp() {
	fs := setupFlags(&options{})
	fs.Usage()
}

//  Hook up the commandLineHelpUsage and commandLineHelpFooter strings
//  to the standard Go flag.Usage function.
func setupUsage(fs *flag.FlagSet) {
	printNonEmpty := func(s string) {
		if s != "" {
			fmt.Fprintf(os.Stderr, "%s\n", s)
		}
	}
	fs.Usage = func() {
		printNonEmpty(CommandLineHelpUsage)
		fs.PrintDefaults()
		printNonEmpty(CommandLineHelpFooter)
	}
}

//  Parse the command line options, validate them, and process them
//  further (e.g. Initialize more complex structs) if need be.
func parseFlags() options {
	var opt options
	fs := setupFlags(&opt)
	fs.Parse(os.Args[1:])
	verifyFlags(&opt, fs)
	// Process the verified options...
	return opt
}
