# golist

[![Go Reference](https://pkg.go.dev/badge/github.com/a-poor/golist.svg)](https://pkg.go.dev/github.com/a-poor/golist)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/a-poor/golist/Go?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-poor/golist?style=flat-square)
![GitHub](https://img.shields.io/github/license/a-poor/golist?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-poor/golist)](https://goreportcard.com/report/github.com/a-poor/golist)
[![Sourcegraph](https://sourcegraph.com/github.com/a-poor/golist/-/badge.svg)](https://sourcegraph.com/github.com/a-poor/golist?badge)

_created by Austin Poor_

A terminal task-list tool for Go. Inspired by the Node package [listr](https://www.npmjs.com/package/listr).

## Example

Here's a quick example of `golist` in action!

Start by creating a list.

```go
list := List{}
```

And let's create some tasks to add to the list. 

`Message` is the text to be displayed and `Action` is the function to be run.

```go
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
```

Next, we'll create a group of 3 sub-tasks.

We'll set the `FailOnError` parameter, so if any of the tasks return an error, the rest will be skipped.

```go
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
```

And let's add that task group and then another task for good measure.

```go
list.AddTask(&tg)
list.AddTask(&Task{
    Message: "Task 4",
    Action: func() error {
        time.Sleep(time.Second / 3)
        return nil
    },
})
```

The `Start` function will start to display the task list and the task statuses.

```go
list.Start()
```

Then, the `Run` function will start to run the tasks syncronously and update the statuses as they complete.

```go
list.Run()
```

And once we're done, we can call `Stop` to stop updating the task status list.

```go
list.Stop()
```

And here's what that looks like in action:

![Sample GIF](https://raw.githubusercontent.com/a-poor/golist/main/etc/sample.gif)

