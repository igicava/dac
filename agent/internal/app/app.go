package app

import (
	"context"
	"encoding/json"
	"log"
	"time"

	pb "dac/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	// установим соединение
	conn, _ := grpc.Dial("orchestrator:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// закроем соединение, когда выйдем из функции
	defer conn.Close()
	grpcClient := pb.NewCalcServiceClient(conn)
	task, err := grpcClient.GETtask(context.TODO(), &pb.GETRequest{})

	if err != nil {
		log.Println("ERROR CONNECT")
	}

	var result struct {
		Task Task `json:"task"`
	}

	err = json.Unmarshal([]byte(task.Result), &result.Task)
	if err != nil {
		log.Println("Error Unmarshal task on agent")
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
	// установим соединение
	conn, _ := grpc.Dial("orchestrator:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// закроем соединение, когда выйдем из функции
	defer conn.Close()
	grpcClient := pb.NewCalcServiceClient(conn)
	_, err = grpcClient.POSTtask(context.TODO(), &pb.POSTRequest{JsonTASK: jsonData})

	if err != nil {
		log.Println("ERROR CONNECT")
	}

}
