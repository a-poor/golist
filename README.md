# golist

![quick & early example](docs/assets/demo.gif)

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/a-poor/golist?label=Version&style=flat-square)](https://pkg.go.dev/github.com/a-poor/golist)
[![Go Reference](https://pkg.go.dev/badge/github.com/a-poor/golist.svg)](https://pkg.go.dev/github.com/a-poor/golist)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/a-poor/golist/Go?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-poor/golist?style=flat-square)
![GitHub](https://img.shields.io/github/license/a-poor/golist?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-poor/golist)](https://goreportcard.com/report/github.com/a-poor/golist)
[![Sourcegraph](https://sourcegraph.com/github.com/a-poor/golist/-/badge.svg)](https://sourcegraph.com/github.com/a-poor/golist?badge)

_created by Austin Poor_

A terminal task-list tool for Go. Inspired by the Node package [listr](https://www.npmjs.com/package/listr) and the [AWS Copilot CLI](https://github.com/aws/copilot-cli).

Check out the documentation [here](https://a-poor.github.io/golist)! _NOTE: The docs site is slightly out of date since some cool new features have been added! Check out this README for updated info and then check back with the docs later!_

In a rush? Check out the [Example](##Example) and the [FAQs](##FAQ)

## Features
* Multi-line updating lists print to the console
* Status updates live (with spinners while processing)
* Nested task groups
* Check if tasks should be skipped or should fail
* Safely print to stdout while the list is being displayed
* Update the task's message while running
* Truncate text output
* Optionally expand/collapse a task-group's subtasks when not running
* Task-groups can skip remaining tasks if failed

## Installation

```sh
go get github.com/a-poor/golist
```

## Dependencies

Just the standard library!

## Example

Here's a quick example of `golist` in action:

```go
// Create a list
list := golist.NewDefaultList()

// Add some tasks!
// This task runs and succeeds
list.AddTask(&golist.Task{
    Message: "Start with this",
    Action: func(c golist.TaskContext) error {
        time.Sleep(time.Second / 2)
        return nil
    },
})
// This task is skipped
list.AddTask(&golist.Task{
    Message: "Then skip this",
    Skip: func(c golist.TaskContext) bool {
        return true
    },
    Action: func(c golist.TaskContext) error {
        time.Sleep(time.Second / 4)
        return nil
    },
})
// And this task runs but fails
list.AddTask(&golist.Task{
    Message: "And finally, this should fail",
    Action: func(c golist.TaskContext) error {
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

And here's a slightly longer example, showing a few more of `golist`'s features:

```go
// Create a new list
list := golist.NewDefaultList()

// Add some tasks
list.AddTask(&golist.Task{
    Message: "Start with this",
    Action: func(c golist.TaskContext) error {
        time.Sleep(time.Second / 2)
        
        // Safely print text to the terminal, even while the
        // List is running
        c.Println("Printed from within a task! 1")
        c.Println("Printed from within a task! 2")
        c.Println("Printed from within a task! 3")
        return nil
    },
})

// Add a task that fails (aka returns an error)
list.AddTask(&golist.Task{
    Message: "This should fail",
    Action: func(c golist.TaskContext) error {
        time.Sleep(time.Second / 3)
        return errors.New("oops")
    },
})

// Add a task that is manually skipped
// (So the `Action` isn't run)
list.AddTask(&golist.Task{
    Message: "Then skip this",
    Skip: func(c golist.TaskContext) bool {
        return true
    },
    Action: func(c golist.TaskContext) error {
        time.Sleep(time.Second / 4)
        c.SetMessage("I renamed myself!")
        time.Sleep(time.Second / 4)
        return nil
    },
})

// Add a group of tasks that will fail if one of
// the tasks returns an error
g := &golist.TaskGroup{
    Message:     "Here's a group",
    FailOnError: true,
}
g.AddTask(&golist.Task{
    Message: "Subtask 1",
    Action: func(c golist.TaskContext) error {
        time.Sleep(time.Second / 2)
        return nil
    },
})

// This task fails causing the rest of the tasks
// in the group to be skipped
g.AddTask(&golist.Task{
    Message: "I'm going to fail",
    Action: func(c golist.TaskContext) error {
        time.Sleep(time.Second / 2)
        return errors.New("womp womp")
    },
})
g.AddTask(&golist.Task{
    Message: "Last one failed so I won't run.",
    Action: func(c golist.TaskContext) error {
        time.Sleep(time.Second / 2)
        return nil
    },
})
list.AddTask(g)

// Add a task that uses the `TaskContext` to update
// it's message while it's running
list.AddTask(&golist.Task{
    Message: "I don't know what to do!",
    Action: func(c golist.TaskContext) error {
        time.Sleep(time.Second / 3)
        c.SetMessage("I renamed myself!")
        time.Sleep(time.Second * 2 / 3)
        return nil
    },
})

// Start to display the task-list
list.Start()

// Start running the actions
list.Run()

// Stop displaying the task-list
list.Stop()
```

## FAQ

* Multi-line updating lists print to the console
* Status updates live (with spinners while processing)
* Nested task groups
* Check if tasks should be skipped or should fail
* Safely print to stdout while the list is being displayed
* Update the task's message while running
* Truncate text output
* Optionally expand/collapse a task-group's subtasks when not running
* Task-groups can skip remaining tasks if failed

_How do I create a new list?_

Most of the time, your best option is to use `golist.NewDefaultList()` -- it will set you up with some sensible defaults. Otherwise, the you'll need to at least set `Writer`, `Delay` and `StatusIndicator`.

_How do I create a new `Task`/`TaskGroup`?_

Both `Task`s and `TaskGroup`s can be created either with their factory functions (`NewTask`, `NewTaskGroup`) or from scratch (`Task{...}`, `TaskGroup{...}`).

## License

[MIT](./LICENSE)

## Contributing

Pull requests are super welcome! For major changes, please open an issue first to discuss what you would like to change. And please make sure to update tests as appropriate.

Or... feel free to just open an issue with some thoughts or suggestions or even just to say Hi and tell me if this library has been helpful!

