package task

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type STATUS string

const (
	StatusWaiting STATUS = "waiting"
	StatusDone    STATUS = "done"
)

func (s STATUS) String() string {
	return string(s)
}

var taskList sync.Map

type Task struct {
	Id       string
	Name     string
	Status   STATUS
	taskFunc func(ctx context.Context) error
}

func New(name string, f func(ctx context.Context) error) (*Task, error) {
	newTask := &Task{
		Id:       fmt.Sprintf("%s%s", name, uuid.New().String()),
		Name:     name,
		Status:   StatusWaiting,
		taskFunc: f,
	}

	addTask(newTask)

	go newTask.call(context.TODO())

	return newTask, nil
}

func GetTaskStatus(id string) string {
	v, loaded := taskList.Load(id)
	if !loaded {
		return ""
	}
	t := v.(*Task)
	return t.Status.String()
}

func addTask(t *Task) {
	// add taskList
	taskList.Store(t.Id, t)
}

//func deleteTask(id string) {
//	taskList.Delete(id)
//}

func (t *Task) call(ctx context.Context) {
	defer func() {
		t.setTaskStatus(StatusDone)
	}()
	err := t.taskFunc(ctx)
	if err != nil {
		return
	}
}

func (t *Task) setTaskStatus(status STATUS) {
	t.Status = status
}

func init() {
	taskList = sync.Map{}
}
