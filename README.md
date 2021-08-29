# golist

_created by Austin Poor_

![quick & early example](./etc/sample.gif)

[![Go Reference](https://pkg.go.dev/badge/github.com/a-poor/golist.svg)](https://pkg.go.dev/github.com/a-poor/golist)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/a-poor/golist/Go?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-poor/golist?style=flat-square)
![GitHub](https://img.shields.io/github/license/a-poor/golist?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-poor/golist)](https://goreportcard.com/report/github.com/a-poor/golist)
[![Sourcegraph](https://sourcegraph.com/github.com/a-poor/golist/-/badge.svg)](https://sourcegraph.com/github.com/a-poor/golist?badge)

A terminal task-list tool for Go. Inspired by the Node package [listr](https://www.npmjs.com/package/listr) and the [AWS Copilot CLI](https://github.com/aws/copilot-cli).

Check out the documentation [here](https://a-poor.github.io/golist)!


## Features
* Multi-line lists print to the console
* Output runs from a gorouting while the user's tasks run
* Status updates live (with spinners while processing)
* Tasks can be skipped or fail

## To Do
* Clean out the public facing API
    * Be consistent with exported vs unexported values
    * Add helper functions & config structs for quickly creating objects
* Add nested lists
    * Each task has a list of optional sub-tasks
    * Or should a `TaskRunner` be an interface? With `Tasks` and `TaskGroups` both implementing the interface? (Should task-groups have actions?)
    * Leave task results behind?
* Pass a context to the action/skip functions
    * Store values for subsequent tasks
    * Values can help later tasks decide if they should stop
* Give the user a way to update the list item text/status/etc. while working
* Run tasks in a group concurrently?
* Fix issue with strings that are too long for one line
    * Truncate text based on terminal width?
    * Calculate & account for strings expected to be multi-line
* Add tests!
    * How do you test modifying terminal output?




