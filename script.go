// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    script.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sat Nov  5 19:59:33 PDT 2011
 *  Description: 
 */

import (
	"exec"
	"io"
)

var (
	Bash   = NewScriptor("bash", "-c", nil)
	Ruby   = NewScriptor("ruby", "-e", nil)
	Perl   = NewScriptor("perl", "-e", nil)
	Python = NewScriptor("python", "-c", nil)
)

type Script interface {
	Execute() error
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
}

type scriptCmd struct {
	*exec.Cmd
}

func (s scriptCmd) SetStdin(r io.Reader)   { s.Stdin = r }
func (s scriptCmd) SetStdout(wr io.Writer) { s.Stdout = wr }
func (s scriptCmd) SetStderr(wr io.Writer) { s.Stderr = wr }

func Pipes(s Script, in io.Reader, out, err io.Writer) {
	s.SetStdin(in)
	s.SetStdout(out)
	s.SetStderr(err)
}

func CombineOutput(s Script, wr io.Writer) io.Writer {
	s.SetStdout(wr)
	s.SetStderr(wr)
	return wr
}

func (s scriptCmd) Execute() error {
	return s.Run()
}

type Scriptor interface {
	NewScript(string, ...string) Script
}

// 	A script command that runs `Name [Flags] Flag Script`
type scriptCtor struct {
	Name  string   // The path (name of the scripting language) to execute.
	Flags []string // Other flags use with the command.
	Flag  string   // The flag 
}

func (sctor *scriptCtor) NewScript(script string, args ...string) Script {
	cargs := make([]string, 2+len(sctor.Flags)+len(args))
	i := 0
	n := copy(cargs[i:], sctor.Flags)
	i += n
	cargs[i] = sctor.Flag
	i++
	cargs[i] = script
	i++
	n = copy(cargs[i:], args)
	i += n
	if n != len(args) {
		panic("bad sizing (<)")
	}
	if i != len(cargs) {
		panic("bad sizing (>)")
	}
	return scriptCmd{exec.Command(sctor.Name, cargs...)}
}

//	Pass NewScriptor the name of a scripting language (executable or path) along
//	with a command line flag that makes the executable execute the contents of a
//	string. If any other command line flags are desired, pass
func NewScriptor(name string, flag string, other []string) Scriptor {
	return &scriptCtor{Name: name, Flag: flag, Flags: other}
}
