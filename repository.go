// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    repository.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sat Nov  5 19:48:23 PDT 2011
 *  Description: 
 */

import ()

type Repository interface {
	Type() string
	Root() string
	Name() (string, error)
	Tags() ([]string, error)
	TagsFetch() error
	TagsPush() error
	TagDelete(string) error
	TagNew(string, ...string) error
}
