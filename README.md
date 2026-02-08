<p align="center">
  <img width="300" src="assets/logo.png" alt="GoZork Logo" />
</p>
<p align="center">A Zork I game ported from <a href="https://github.com/historicalsource/zork1" target="_blank">ZIL source</a> to Golang.</p>
<p align="center">
  <a href="https://goreportcard.com/report/github.com/ajdnik/gozork" target="_blank"><img src="https://goreportcard.com/badge/github.com/ajdnik/gozork" alt="Go Report" /></a>
</p>

## Description

Zork I is a 1980 interactive fiction game written by Marc Blank, Dave Lebling, Bruce Daniels and Tim Anderson and published by Infocom. To learn more about the history of the game feel free to read [Zork I: The Great Underground Empire](https://medium.com/@r.ajdnik/zork-the-great-underground-empire-cda94623861c) on Medium.

## Prerequisites

- **Go 1.25+** &mdash; <https://go.dev/dl/>
- **golangci-lint v2** (optional, for linting) &mdash; <https://golangci-lint.run/welcome/install/>

## Quick start

```bash
git clone https://github.com/ajdnik/gozork.git
cd gozork
make run
```

## Makefile targets

| Target       | Description                                       |
|--------------|---------------------------------------------------|
| `make build` | Compile the `gozork` binary                       |
| `make run`   | Build and launch the game                          |
| `make test`  | Run all tests (verbose, no cache)                  |
| `make vet`   | Run `go vet` static analysis                       |
| `make lint`  | Run `golangci-lint v2` (requires installation)     |
| `make fmt`   | Format all Go source files with `gofmt`            |
| `make cover` | Run tests and print per-package coverage summary   |
| `make check` | Run fmt, vet, lint, and test in sequence            |
| `make clean` | Remove build artifacts                             |

## Project structure

```
gozork/
  engine/   # Reusable text-adventure engine (parser, objects, clock, I/O)
  game/     # Zork I game content (rooms, items, NPCs, action handlers)
  main.go   # Entry point
```

## Usage

```
$ ./gozork
ZORK I: The Great Underground Empire
Infocom interactive fiction - a fantasy story
Copyright (c) 1981, 1982, 1983, 1984, 1985, 1986 Infocom, Inc. All rights reserved.
ZORK is a registered trademark of Infocom, Inc.
Release 119 / Serial number 880429

West of House
You are standing in an open field west of a white house, with a boarded front door.
There is a small mailbox here.

>
```
