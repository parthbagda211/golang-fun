package main

import (
	"flag"
	"fmt"
	"os"
	"cli-project-go/utils"
)

func main() {
	add := flag.String("add", "", "Add a task")
	list := flag.Bool("list", false, "List all tasks")
	done := flag.Int("done", 0, "Mark task as done")
	delete := flag.Int("delete", 0, "Delete task")
	flag.Parse()

	tasks,err := utils.LoadTasks()
	if err != nil {
		fmt.Println("Error loading tasks", err)
		os.Exit(1)
	}

	switch {
	case *add != "":
		tasks.AddTask(*add)
		utils.SaveTasks(tasks)
		fmt.Println("Task added")
	case *list:
		tasks.ListTasks()
	case *done != 0:
		
		err:= tasks.MarkTaskAsDone(*done)
		if err != nil {
			fmt.Println("Error marking task as done", err)
	
		} else {
			utils.SaveTasks(tasks)
			fmt.Println("Task marked as done")
		}

	case *delete != 0:
		err:= tasks.DeleteTask(*delete)
		if err != nil {
			fmt.Println("Error deleting task", err)
		} else {
			utils.SaveTasks(tasks)
			fmt.Println("Task deleted")
		}
	}
}