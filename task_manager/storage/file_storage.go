package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"task_manager/task"
)

type FileStorage struct {
	FileName string
}

func (fs FileStorage) Load() ([]task.Task, error) {
	file, err := os.Open(fs.FileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []task.Task{}, nil
		}
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []task.Task{}, nil
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, errors.New("corrupted JSON file")
	}

	return tasks, nil
}

func (fs FileStorage) Save(tasks []task.Task) error {
	file, err := os.Create(fs.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	return enc.Encode(tasks)
}
