package utils

import (
	"encoding/json"
	"os"

)

const FilePaht = "tasks.json"

func LoadTasks() (*TaskList, error) {
	file, err := os.Open(FilePaht)
	if err != nil {
		if os.IsNotExist(err) {
			return &TaskList{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var taskList TaskList
	err = json.NewDecoder(file).Decode(&taskList)
	if err != nil {
		return nil, err
	}
	return &taskList, nil
}

func SaveTasks(t *TaskList) error {
	file, err := os.Create(FilePaht)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(t)
	if err != nil {
		return err
	}
	return nil
}
