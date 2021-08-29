# golist

[![Go Reference](https://pkg.go.dev/badge/github.com/a-poor/golist.svg)](https://pkg.go.dev/github.com/a-poor/golist)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/a-poor/golist/Go?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-poor/golist?style=flat-square)
![GitHub](https://img.shields.io/github/license/a-poor/golist?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-poor/golist)](https://goreportcard.com/report/github.com/a-poor/golist)
[![Sourcegraph](https://sourcegraph.com/github.com/a-poor/golist/-/badge.svg)](https://sourcegraph.com/github.com/a-poor/golist?badge)

_created by Austin Poor_

A terminal task-list tool for Go. Inspired by the Node package [listr](https://www.npmjs.com/package/listr) and the [AWS Copilot CLI](https://github.com/aws/copilot-cli).

## Quick Example

```go
// Create a List
list := List{}

// Add some tasks
list.AddTask(&Task{
    Message: "Task 1",
    Action: func() error {
        time.Sleep(time.Second / 2)
        return nil
    },
})
list.AddTask(&Task{
    Message: "Task 2",
    Action: func() error {
        time.Sleep(time.Second / 4)
        return nil
    },
})

// Create a group of tasks
tg := TaskGroup{
    Message:     "Task Group 3",
    FailOnError: true,
}
tg.AddTask(&Task{
    Message: "Task 3a",
    Action: func() error {
        time.Sleep(time.Second / 4)
        return nil
    },
})
tg.AddTask(&Task{
    Message: "Task 3b",
    Action: func() error {
        time.Sleep(time.Second / 2)
        return errors.New("oh no")
    },
})
tg.AddTask(&Task{
    Message: "Task 3c",
    Action: func() error {
        time.Sleep(time.Second / 3)
        return nil
    },
})

// Add the TaskGroup
list.AddTask(&tg)
list.AddTask(&Task{
    Message: "Task 4",
    Action: func() error {
        time.Sleep(time.Second / 3)
        return nil
    },
})

// Start displaying the task list & statuses
fmt.Println("Starting...")
list.Start()

// Run the tasks (syncronously)
list.Run()

// Stop displaying the task list & statuses
list.Stop()
fmt.Println("Done.")
```

And the result is...

![Sample GIF](https://raw.githubusercontent.com/a-poor/golist/main/etc/sample.gif)

