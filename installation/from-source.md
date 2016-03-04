---
layout: default
title: Installation From Source
description: How to install Inbucket from source
---

## Prerequisites

You will need functioning [Go][1], [Git][2], and [Mercurial][3] installations
for this to work.

If you are running on Windows, you will need to perform the steps in this
document from within Git Bash shell, otherwise Go will be unable to clone
(download) the Inbucket source code and its dependencies.

If you don't have a `GOPATH` setup, create a directory for holding Go source
code and export the `GOPATH` variable:

    mkdir ~/gocode
    export GOPATH=~/gocode

## Building

Grab the Inbucket source code and compile the daemon:

    go get -v github.com/jhillyerd/inbucket

## Launching

Test launching the daemon on Linux & OS X:

    cd $GOPATH/src/github.com/jhillyerd/inbucket
    $GOPATH/bin/inbucket etc/inbucket.conf

Test launching the daemon on Windows:

    cd $GOPATH/src/github.com/jhillyerd/inbucket
    $GOPATH/bin/inbucket etc/win-sample.conf

By default the SMTP server will be listening on localhost port 2500 and the web
interface will be available at [localhost:9000](http://localhost:9000/).

## Next Steps

Please see our [RedHat](redhat.html) and [Ubuntu](ubuntu.html) installation
guides for information on configuring Inbucket to run as a daemon.

[1]: http://golang.org/
[2]: http://git-scm.com/
[3]: http://mercurial.selenic.com/