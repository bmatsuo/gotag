About Gotag
=============

Gotag is used to keep Go version tags up to date in your repositories.
Go version tags are used by [Goinstall] [] to distinguish which version of your
project to install.

    ```sh
    git tag -a -m "Latest build for Go version rXX.X" go.rXX.X
    ```

This is great and all. Things were crappy before. But, updating tags (with git)
is a little more of a pain than I would like. And, the annotation message for
the tag is always the same.

    ```sh
    git tag -d go.rXX.X
    git tag -a -m "Latest build for Go version rXX.X" go.rXX.X
    ```

So Gotag was created to solve this little annoyance in life.

    ```
    $ gotag -f
    testing build: PASS
    found tag go.rXX.X
    deleting go.rXX.X
    creating tag go.rXX.X "Latest build for Go version rXX.X"
    $
    ```

Note: That doesn't happen yet. I'm still working on it.

Documentation
=============

Usage
-----

Run gotag with the command

    gotag [-f]

Prerequisites
-------------

[Install Go] []

Installation
-------------

Use goinstall to install gotag

    goinstall github.com/bmatsuo/gotag

General Documentation
---------------------

Use godoc to vew the documentation for gotag

    godoc github.com/bmatsuo/gotag

Or alternatively, use a godoc http server

    godoc -http=:6060

and view the [Godoc URL][]

Author
======

Bryan Matsuo <bryan.matsuo@gmail.com>

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.

Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.

[goisntall]: http://golang.org/cmd/goinstall "Goinstall"
[install go]: http://golang.org/doc/install.html "Install Go"
[godoc url]: http://localhost:6060/pkg/github.com/bmatsuo/go-script/script/ "Godoc URL"
