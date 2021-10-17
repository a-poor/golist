package golist

// TaskGroup represents a group of TaskRunners
// for running nested tasks within a TaskList
type TaskGroup struct {
	Message                 string                 // The message to be displayed
	Tasks                   []TaskRunner           // A list of tasks to run
	Skip                    func(TaskContext) bool // Is run before the task starts. If returns true, the task isn't run
	FailOnError             bool                   // If true, the task group stops on the first error
	HideTasksWhenNotRunning bool                   // If true, the task group only show its sub-task-runners when actively running

	err    error      // The error that occurred during the last task
	status TaskStatus // The status of the task
}

// NewTaskGroup creates a new default TaskGroup
func NewTaskGroup() *TaskGroup {
	return &TaskGroup{}
}

// AddTask adds a TaskRunner to this TaskGroup's tasks
func (tg *TaskGroup) AddTask(t TaskRunner) {
	tg.Tasks = append(tg.Tasks, t)
}

// Run runs the TaskRunners in this TaskGroup
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

// SetMessage sets the display message for this TaskGroup
func (tg *TaskGroup) SetMessage(m string) {
	tg.Message = m
}

// GetError returns this TaskGroup's error, if any
func (tg *TaskGroup) GetError() error {
	return tg.err
}

// SetError sets this TaskGroup's error value
func (tg *TaskGroup) SetError(err error) {
	tg.err = err
}

// GetStatus returns this TaskGroup's TaskStatus
func (tg *TaskGroup) GetStatus() TaskStatus {
	return tg.status
}

// GetStatus sets this TaskGroup's TaskStatus
func (tg *TaskGroup) SetStatus(s TaskStatus) {
	tg.status = s
}

// GetTaskStates returns a slice of TaskStates for this TaskGroup
// representing it's current state as well as the state of its sub-tasks
// (by calling GetTaskStates on each of its sub-tasks). TaskStates store
// a TaskRunners message, status, and tree-depth, and are passed up to
// the parent List for printing.
func (tg *TaskGroup) GetTaskStates() []*TaskState {
	messages := []*TaskState{{
		Status:  tg.GetStatus(),
		Message: tg.Message,
	}}
	if !tg.HideTasksWhenNotRunning || tg.GetStatus() == TaskInProgress {
		for _, t := range tg.Tasks {
			msgs := t.GetTaskStates()
			for _, m := range msgs {
				m.Depth++
			}
			messages = append(messages, msgs...)
		}
	}
	return messages
}
