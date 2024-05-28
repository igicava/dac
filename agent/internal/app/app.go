package app

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Структура таска
type Task struct {
    ID            string  `json:"id"`
    Arg1          float64 `json:"arg1"`
    Arg2          float64 `json:"arg2"`
    Operation     string  `json:"operation"`
    OperationTime int     `json:"operation_time"`
}

// Этот бро получает выражения
func GetTask() *Task {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var result struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil
	}
	return &result.Task
}

// Этот бро считает выражения
func ComputeTask(task Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			log.Println("Division by zero")
			return 0
		}
		return task.Arg1 / task.Arg2
	}
	log.Println("Unknown operation:", task.Operation)
	return 0
}

// А этот бро отправляет результат оркестратору
func SendResult(id string, result float64) {
	taskResult := struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}{
		ID:     id,
		Result: result,
	}

	jsonData, err := json.Marshal(taskResult)
	if err != nil {
		log.Println("Error marshalling result:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending result:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error from server:", resp.Status)
	}
}
