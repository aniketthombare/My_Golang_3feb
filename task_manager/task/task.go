package task

import (
	"fmt"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

func New(id int, title string) Task {
	return Task{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
}

func (t *Task) Complete() {
	t.Completed = true
}

func (t Task) String() string {
	status := " Pending"
	if t.Completed {
		status = " Completed"
	}

	return fmt.Sprintf(
		"[%d] %s | %s | Created: %s",
		t.ID,
		t.Title,
		status,
		t.CreatedAt.Format(time.RFC822),
	)
}
