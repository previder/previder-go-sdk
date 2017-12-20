package client

import (
	"errors"
	"time"
)

const taskBasePath = "configurationitem/task/"

type TaskService interface {
	List() (*[]Task, error)
	Get(id string) (*Task, error)
	WaitFor(id string, timeoutDuration time.Duration) (*Task, error)
	WaitForTask(task *Task, timeoutDuration time.Duration) (*Task, error)
}

type TaskServiceOp struct {
	client *BaseClient
}

type Task struct {
	Id                string
	TaskDate          *uint64
	StartedDate       *uint64
	CompletedDate     *uint64
	Success           bool
	Error             string
	ConfigurationItem ConfigurationItem
}

type ConfigurationItem struct {
	Id   string
	Name string
}

func (c *TaskServiceOp) List() (*[]Task, error) {
	task := new([]Task)
	err := c.client.Get(taskBasePath, task)
	return task, err
}

func (c *TaskServiceOp) Get(id string) (*Task, error) {
	task := new(Task)
	err := c.client.Get(taskBasePath+id, task)
	return task, err
}

func (c *TaskServiceOp) WaitForTask(task *Task, timeoutDuration time.Duration) (*Task, error) {
	return c.WaitFor(task.Id, timeoutDuration)
}

func (c *TaskServiceOp) WaitFor(id string, timeoutDuration time.Duration) (*Task, error) {
	timeout := time.After(timeoutDuration)
	tick := time.Tick(3 * time.Second)
	for {
		select {
		case <-timeout:
			return nil, errors.New("timed out")
		case <-tick:
			task, err := c.Get(id)
			if err != nil {
				return nil, err
			}
			if task.CompletedDate != nil {
				if task.Success {
					return task, nil
				} else {
					return task, errors.New(task.Error)
				}
			}
		}
	}
}
