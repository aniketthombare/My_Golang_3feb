package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aniketthombare/My_Go_1/task_manager/manager"
	"github.com/aniketthombare/My_Go_1/task_manager/storage"
)

const (
	AddTaskOption      = 1
	ListTasksOption    = 2
	ViewTaskOption     = 3
	CompleteTaskOption = 4
	ExitOption         = 5
)

func main() {
	fs := storage.FileStorage{FileName: "tasks.json"}

	mgr, err := manager.New(fs)
	if err != nil {
		fmt.Println("Failed to load tasks:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()

		fmt.Print("Select option: ")
		choiceStr, _ := reader.ReadString('\n')
		choice, _ := strconv.Atoi(strings.TrimSpace(choiceStr))

		switch choice {
		case AddTaskOption:
			handleAddTask(reader, mgr)

		case ListTasksOption:
			mgr.PrintTasks()

		case ViewTaskOption:
			handleViewTask(reader, mgr)

		case CompleteTaskOption:
			handleCompleteTask(reader, mgr)

		case ExitOption:
			fmt.Println("Saving & exiting...")

			resultChan := make(chan error)
			mgr.SaveAsync(resultChan)

			select {
			case err := <-resultChan:
				if err != nil {
					fmt.Println("Save error:", err)
				} else {
					fmt.Println("Tasks saved successfully.")
				}
			case <-time.After(3 * time.Second):
				fmt.Println("Save timed out!")
			}

			return

		default:
			fmt.Println("Invalid option")
		}
	}
}

func printMenu() {
	fmt.Println(`
==== TASK MANAGER ====
1. Add Task
2. List Tasks
3. View Task by ID
4. Complete Task
5. Exit
`)
}

func handleAddTask(reader *bufio.Reader, mgr *manager.Manager) {
	fmt.Print("Enter ID: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	fmt.Print("Enter title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	if title == "" {
		fmt.Println("Title cannot be empty")
		return
	}

	if err := mgr.AddTask(id, title); err != nil {
		fmt.Println("Error:", err)
		return
	}

	saveAsync(mgr)
}

func handleViewTask(reader *bufio.Reader, mgr *manager.Manager) {
	fmt.Print("Enter ID: ")
	idStr, _ := reader.ReadString('\n')
	id, _ := strconv.Atoi(strings.TrimSpace(idStr))

	task, err := mgr.GetTask(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(task.String())
}

func handleCompleteTask(reader *bufio.Reader, mgr *manager.Manager) {
	fmt.Print("Enter ID: ")
	idStr, _ := reader.ReadString('\n')
	id, _ := strconv.Atoi(strings.TrimSpace(idStr))

	if err := mgr.CompleteTask(id); err != nil {
		fmt.Println(err)
		return
	}

	saveAsync(mgr)
}

func saveAsync(mgr *manager.Manager) {
	resultChan := make(chan error)
	mgr.SaveAsync(resultChan)

	go func() {
		select {
		case err := <-resultChan:
			if err != nil {
				fmt.Println("Auto-save failed:", err)
			} else {
				fmt.Println("Auto-saved")
			}
		case <-time.After(2 * time.Second):
			fmt.Println("Auto-save timeout")
		}
	}()
}
