package manager

import (
	"errors"
	"fmt"
	"sync"

	"task_manager/storage"
	"task_manager/task"
)

type Manager struct {
	storage storage.Storage

	tasks   []task.Task
	taskMap map[int]*task.Task

	mu sync.Mutex
}

func New(storage storage.Storage) (*Manager, error) {
	tasks, err := storage.Load()
	if err != nil {
		return nil, err
	}

	m := &Manager{
		storage: storage,
		tasks:   tasks,
		taskMap: make(map[int]*task.Task),
	}

	for i := range tasks {
		t := &m.tasks[i]
		m.taskMap[t.ID] = t
	}

	return m, nil
}

func (m *Manager) AddTask(id int, title string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.taskMap[id]; exists {
		return errors.New("task with this ID already exists")
	}

	t := task.New(id, title)
	m.tasks = append(m.tasks, t)
	m.taskMap[id] = &m.tasks[len(m.tasks)-1]

	return nil
}

func (m *Manager) ListTasks() []task.Task {
	return m.tasks
}

func (m *Manager) GetTask(id int) (*task.Task, error) {
	if t, ok := m.taskMap[id]; ok {
		return t, nil
	}
	return nil, errors.New("task not found")
}

func (m *Manager) CompleteTask(id int) error {
	t, err := m.GetTask(id)
	if err != nil {
		return err
	}

	t.Complete()
	return nil
}

func (m *Manager) SaveAsync(resultChan chan error) {
	go func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		err := m.storage.Save(m.tasks)
		resultChan <- err
	}()
}

func (m *Manager) PrintTasks() {
	fmt.Println("\n----- TASK LIST -----")
	for _, t := range m.tasks {
		fmt.Println(t.String())
	}
	fmt.Println("---------------------")
}
