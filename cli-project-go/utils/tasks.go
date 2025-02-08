package utils

import (
    "errors"
	"fmt"
)

type Task struct {
	ID int
	Description string
	Done bool
}

type TaskList struct {
	Tasks []Task
	NextId int
}

func (t *TaskList) AddTask(description string) {
	task := Task{t.NextId, description, false}
	t.Tasks = append(t.Tasks, task)
	t.NextId++
}

func (t *TaskList) ListTasks(){ 
	if len(t.Tasks) == 0 {
		fmt.Println("No tasks")
		return
	}

	for _,task := range t.Tasks{
		status := "pending"
		if task.Done {
			status = "done"
		}
		fmt.Printf("%d: %s (%s)\n", task.ID, task.Description, status)
	}
}

func (t *TaskList) MarkTaskAsDone(id int) error {
	for i:= range t.Tasks{
		if t.Tasks[i].ID == id {
			t.Tasks[i].Done = true
			return nil
		}
	}
	return errors.New("task not found")
}

func (t *TaskList) DeleteTask(id int) error {
	for i:= range t.Tasks{
		if t.Tasks[i].ID == id {
			t.Tasks = append(t.Tasks[:i], t.Tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}