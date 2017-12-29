# fscrub
[![Go Report Card](https://goreportcard.com/badge/github.com/playnet-public/fscrub)](https://goreportcard.com/report/github.com/playnet-public/fscrub)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Build Status](https://travis-ci.org/playnet-public/fscrub.svg?branch=master)](https://travis-ci.org/playnet-public/fscrub)
[![Docker Repository on Quay](https://quay.io/repository/playnet/fscrub/status "Docker Repository on Quay")](https://quay.io/repository/playnet/fscrub)
[![Join Discord at https://discord.gg/dWZkR6R](https://img.shields.io/badge/style-join-green.svg?style=flat&label=Discord)](https://discord.gg/dWZkR6R)

Tool watching and crawling the file-system for files containing sensitive information to clean them.

## Description

The whole idea for fscrub arose on a public forum where some people posted log files from their apache or game servers to get help from the community.
Those log files contained not only ip addresses of the clients but also other sensitive information.

To prevent this from happening we had the idea to provide the board with a tool to monitor the servers attachments directory and acting on relevant files.
This project is WIP and only a small side project.
Aside from that, any community project is free to use it so their users and their user's clients stay safe.  

fscrub has several modus operandi:
* watching a provided folder for file-system changes and acting on creation/change of files
* crawling a provided folder for files of interest and acting on find
* watching a provided folder for new changes while also crawling to ensure past changes are covered

The actions fscrub is planed to take are:
* finding personal or security relevant data based on provided patterns and replacing them (ip's, passwords, hostnames, etc.)
* checking files against virus check api's or antivir tools
* replacing found data or files with scrubbed ones or deleting them
* moving found files to a safe backup location
* notifying a pool of users on certain events


## Dependencies

This project has a pretty complex Makefile and therefore requires `make`.

Go Version: 1.8

Install all further requirements by running `make deps`

## Usage

```
fscrub -watch -crawl -dir=/var/www/attachments -bckdir=/var/www/backups/fscrub -patterns=/var/opt/fscrub/patterns.json
```

It is possible to provide multiple dirs to handle. To do so, simply use the `-dir` parameter multiple times:
```
fscrub -dir=./pkg -dir=./cmd
```

## Development

This project is using a [basic template](github.com/playnet-public/gocmd-template) for developing PlayNet command-line tools. Refer to this template for further information and usage docs.
The Makefile is configurable to some extent by providing variables at the top.
Any further changes should be thought of carefully as they might brake CI/CD compatibility.

One project might contain multiple tools whose main packages reside under `cmd`. Other packages like libraries go into the `pkg` directory.
Single projects can be handled by calling `make toolname maketarget` like for example:
```
make template dev
```
All tools at once can be handled by calling `make full maketarget` like for example:
```
make full build
```
Build output is being sent to `./build/`.

If you only package one tool this might seam slightly redundant but this is meant to provide consistence over all projects.
To simplify this, you can simply call `make maketarget` when only one tool is located beneath `cmd`. If there are more than one, this won't do anything (including not return 1) so be careful.

## Contributions

Pull Requests and Issue Reports are welcome.
If you are interested in contributing, feel free to [get in touch](https://discord.gg/WbrXWJB)