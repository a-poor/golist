package golist

type TaskGroup struct {
	Message     string                 // The message to be displayed
	Tasks       []TaskRunner           // A list of tasks to run
	Skip        func(TaskContext) bool // Is run before the task starts. If returns true, the task isn't run
	FailOnError bool                   // If true, the task group stops on the first error

	err    error      // The error that occurred during the last task
	status TaskStatus // The status of the task
}

func NewTaskGroup() *TaskGroup {
	return &TaskGroup{}
}

func (tg *TaskGroup) AddTask(t TaskRunner) {
	tg.Tasks = append(tg.Tasks, t)
}

func (tg *TaskGroup) Run(parentContext TaskContext) error {
	// Create a context
	c := tg.createContext(parentContext)

	// Check if the task should be skipped
	if tg.Skip != nil && tg.Skip(c) {
		tg.SetStatus(TaskSkipped)
		return nil
	}

	// Run the tasks
	tg.SetStatus(TaskInProgress)
	var err error
	var skipRemaining bool
	for _, t := range tg.Tasks {
		if skipRemaining {
			t.SetStatus(TaskSkipped)
			continue
		}
		err = t.Run(c)
		if err != nil && tg.FailOnError {
			skipRemaining = true
		}
	}

	// Set the final status
	if err != nil {
		tg.SetStatus(TaskFailed)
	} else {
		tg.SetStatus(TaskCompleted)
	}
	return err
}

// createContext creates a TaskContext for the task
func (t *TaskGroup) createContext(parentContext TaskContext) TaskContext {
	return &taskContext{
		setMessage: func(msg string) {
			t.SetMessage(msg)
		},
		println: func(a ...interface{}) error {
			return parentContext.Println(a...)
		},
		printfln: func(f string, a ...interface{}) error {
			return parentContext.Printfln(f, a...)
		},
	}
}

func (tg *TaskGroup) SetMessage(m string) {
	tg.Message = m
}

func (tg *TaskGroup) GetError() error {
	return tg.err
}

func (tg *TaskGroup) SetError(err error) {
	tg.err = err
}

func (tg *TaskGroup) GetStatus() TaskStatus {
	return tg.status
}

func (tg *TaskGroup) SetStatus(s TaskStatus) {
	tg.status = s
}

func (tg *TaskGroup) GetTaskStates() []*TaskState {
	messages := []*TaskState{{
		Status:  tg.GetStatus(),
		Message: tg.Message,
	}}
	for _, t := range tg.Tasks {
		msgs := t.GetTaskStates()
		for _, m := range msgs {
			m.Depth++
		}
		messages = append(messages, msgs...)
	}
	return messages
}
