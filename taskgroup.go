package main

type TaskGroup struct {
	Tasks []Task
}

func (tg *TaskGroup) Run() error {
	return nil
}
