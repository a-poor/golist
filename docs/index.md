# golist

![Sample GIF](assets/sample.gif)

[![Go Reference](https://pkg.go.dev/badge/github.com/a-poor/golist.svg)](https://pkg.go.dev/github.com/a-poor/golist)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/a-poor/golist/Go?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-poor/golist?style=flat-square)
![GitHub](https://img.shields.io/github/license/a-poor/golist?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-poor/golist)](https://goreportcard.com/report/github.com/a-poor/golist)
[![Sourcegraph](https://sourcegraph.com/github.com/a-poor/golist/-/badge.svg)](https://sourcegraph.com/github.com/a-poor/golist?badge)

_created by Austin Poor_

## About

A terminal task-list tool for Go. Inspired by the Node package [listr](https://www.npmjs.com/package/listr).

## Quickstart

Install golist with:

```bash
$ go get github.com/a-poor/golist
```

And create a task-list:

```go
// Create a list
list := golist.List{}

// Add some tasks!
// This task runs and succeeds
list.AddTask(&golist.Task{
    Message: "Start with this",
    Action: func() error {
        time.Sleep(time.Second / 2)
        return nil
    },
})
// This task is skipped
list.AddTask(&golist.Task{
    Message: "Then skip this",
    Skip: func() bool {
        return true
    },
    Action: func() error {
        time.Sleep(time.Second / 4)
        return nil
    },
})
// And this task runs but fails
list.AddTask(&golist.Task{
    Message: "And finally, this should fail",
    Action: func() error {
        time.Sleep(time.Second / 3)
        return errors.New("oops")
    },
})

// Start displaying the task status
list.Start()

// Run the tasks
list.Run()

// Stop displaying the task status
list.Stop()
```

## Etc

Let me know what you think of `golist`! I'd love any feedback you have. 

Please feel free to submit an issue or a pr!


