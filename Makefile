# Modified the basic makefiles referred to from the
# Go home page.
#
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=gotag
GOFILES=\
		project.go\
		git.go\
		repository.go\
        options.go\
        gotag.go\

include $(GOROOT)/src/Make.cmd


test:
	gotest

