# golist

![quick & early example](docs/assets/sample.gif)

[![Go Reference](https://pkg.go.dev/badge/github.com/a-poor/golist.svg)](https://pkg.go.dev/github.com/a-poor/golist)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/a-poor/golist/Go?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-poor/golist?style=flat-square)
![GitHub](https://img.shields.io/github/license/a-poor/golist?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-poor/golist)](https://goreportcard.com/report/github.com/a-poor/golist)
[![Sourcegraph](https://sourcegraph.com/github.com/a-poor/golist/-/badge.svg)](https://sourcegraph.com/github.com/a-poor/golist?badge)

_created by Austin Poor_

A terminal task-list tool for Go. Inspired by the Node package [listr](https://www.npmjs.com/package/listr) and the [AWS Copilot CLI](https://github.com/aws/copilot-cli).

Check out the documentation [here](https://a-poor.github.io/golist)!


## Features
* Multi-line lists print to the console
* Output runs from a gorouting while the user's tasks run
* Status updates live (with spinners while processing)
* Tasks can be skipped or fail

## Installation

```sh
go get github.com/a-poor/golist
```

## Dependencies

...

## Example

Here's a quick example of `golist` in action:

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

## License

[MIT](./LICENSE)

## Contributing

Pull requests are super welcome! For major changes, please open an issue first to discuss what you would like to change. And please make sure to update tests as appropriate.

Or... feel free to just open an issue with some thoughts or suggestions or even just to say Hi and tell me if this library has been helpful!

