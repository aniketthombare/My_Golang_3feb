package storage

import "task_manager/task"

type Storage interface {
	Load() ([]task.Task, error)
	Save([]task.Task) error
}
