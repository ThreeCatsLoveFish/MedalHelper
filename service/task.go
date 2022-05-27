package service

import "sync"

// Task aggregate user info and corresponding action
type Task struct {
	User
	actions []IAction
}

func NewTask(user User, actions []IAction) Task {
	return Task{
		User:    user,
		actions: actions,
	}
}

func (task *Task) Start() {
	wg := sync.WaitGroup{}
	for _, action := range task.actions {
		wg.Add(1)
		go action.Exec(task.User, &wg, action)
	}
	wg.Wait()
}
